package livebili

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kohmebot/livebili/request"
	"github.com/kohmebot/plugin/pkg/chain"
	"github.com/kohmebot/plugin/pkg/gopool"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"io"
	"sync"
	"time"
)

func (b *biliPlugin) doCheckDynamic() error {
	var uids []int64
	for _, uid := range b.conf.Uids {
		uids = append(uids, uid)
	}
	errChan := make(chan error, len(uids))
	defer close(errChan)
	for _, uid := range uids {
		gopool.Go(func() {
			errChan <- b.doCheckOneDynamic(uid)
		})
	}
	var err error
	for i := 0; i < len(uids); i++ {
		doErr := <-errChan
		if doErr != nil {
			err = errors.Join(err, fmt.Errorf("checkDynamicError %w", doErr))
		}
	}
	return nil
}

func (b *biliPlugin) doCheckOneDynamic(uid int64) error {

	resp, err := request.DoGet(fmt.Sprintf("https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/space?host_mid=%d", uid), b.conf.Cookies)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dynamic := DynamicResp{}
	err = json.Unmarshal(body, &dynamic)
	if err != nil {
		return err
	}
	if dynamic.Code != 0 {
		return fmt.Errorf("code: %d,msg: %s", dynamic.Code, dynamic.Message)
	}

	updates, err := b.updateDynamic(uid, &dynamic)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}

	b.env.RangeBot(func(ctx *zero.Ctx) bool {
		b.groups.RangeGroup(func(group int64) bool {
			wg.Add(1)
			gopool.Go(func() {
				defer wg.Done()
				for _, update := range updates {
					b.sendDynamic(ctx, update)
					// 发送每个动态间等2s
					time.Sleep(2 * time.Second)
				}
			})
			return true
		})
		return true
	})
	wg.Wait()
	return nil

}

func (b *biliPlugin) updateDynamic(uid int64, dynamic *DynamicResp) (updates []*Dynamic, err error) {
	if len(dynamic.Data.Items) <= 0 {
		return nil, nil
	}
	db, err := b.env.GetDB()
	if err != nil {
		return nil, err
	}
	lastPubTime := dynamic.Data.Items[0].Modules.PutTs
	record := &DynamicRecord{}
	record.Uid = uid
	err = db.First(&record, uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有记录，插入一条新记录
			record = &DynamicRecord{
				Uid:         uid,
				LastPubTime: lastPubTime,
			}
			err = db.Create(&record).Error
		}
		return nil, err
	}

	// 说明动态未更新
	if record.LastPubTime >= lastPubTime {
		return nil, nil
	}

	// 找出更新的动态
	for _, item := range dynamic.Data.Items {
		if item.Modules.PutTs > record.LastPubTime {
			updates = append(updates, &item)
		}
	}
	return updates, nil

}

func (b *biliPlugin) sendDynamic(ctx *zero.Ctx, dynamic *Dynamic) {
	switch dynamic.Type {
	case "DYNAMIC_TYPE_AV":
		b.onAv(ctx, &dynamic.Modules)
	case "DYNAMIC_TYPE_DRAW":
		b.onDraw(ctx, &dynamic.Modules)
	case "DYNAMIC_TYPE_WORD":
		b.onWord(ctx, &dynamic.Modules)
	default:
		logrus.Warnf("unknown dynamic type: %s", dynamic.Type)
	}
}

// 投稿了视频
func (b *biliPlugin) onAv(ctx *zero.Ctx, dynamic *DynamicModules) {
	userName := dynamic.ModuleAuthor.Name
	pubTime := dynamic.ModuleAuthor.PubTime

	title := dynamic.Archive.Title
	cover := dynamic.Archive.Cover
	url := dynamic.Archive.JumpUrl

	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", userName)),
		message.Text(fmt.Sprintf("%s投稿了视频", pubTime)),
		message.Text(fmt.Sprintf("【%s】", title)),
		message.Image(cover),
		message.Text(url),
	)
	b.groups.RangeGroup(func(group int64) bool {
		gopool.Go(func() {
			ctx.SendGroupMessage(group, msgChain)
		})
		return true
	})

}

// 带图动态
func (b *biliPlugin) onDraw(ctx *zero.Ctx, dynamic *DynamicModules) {
	userName := dynamic.ModuleAuthor.Name
	pubTime := dynamic.ModuleAuthor.PubTime

	text := dynamic.Desc.Text

	var imgMsg []message.MessageSegment
	for _, item := range dynamic.Draw.Items {
		imgMsg = append(imgMsg, message.Image(item.Src))
	}

	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", userName)),
		message.Text(fmt.Sprintf("%s发布了动态", pubTime)),
		message.Text(text),
	)
	if len(imgMsg) > 0 {
		msgChain.Line()
		msgChain.Split(imgMsg...)
	}
	b.groups.RangeGroup(func(group int64) bool {
		gopool.Go(func() {
			ctx.SendGroupMessage(group, msgChain)
		})
		return true
	})

}

// 纯文字动态
func (b *biliPlugin) onWord(ctx *zero.Ctx, dynamic *DynamicModules) {
	userName := dynamic.ModuleAuthor.Name
	pubTime := dynamic.ModuleAuthor.PubTime

	text := dynamic.Desc.Text

	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", userName)),
		message.Text(fmt.Sprintf("%s发布了动态", pubTime)),
		message.Text(text),
	)
	b.groups.RangeGroup(func(group int64) bool {
		gopool.Go(func() {
			ctx.SendGroupMessage(group, msgChain)
		})
		return true
	})
}
