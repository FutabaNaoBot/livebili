package livebili

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kohmebot/livebili/gopool"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

func (b *biliPlugin) doCheckLive() error {
	var uids []int64
	for _, uid := range b.conf.Uids {
		uids = append(uids, uid)
	}
	data := map[string]interface{}{
		"uids": uids,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post("https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	live := LiveResp{}
	err = json.Unmarshal(body, &live)
	if err != nil {
		return err
	}
	if live.Code != 0 {
		return fmt.Errorf("code: %d,msg: %s", live.Code, live.Msg)
	}

	for _, info := range live.Data {
		err = b.sendRoomInfo(&info)
		if err != nil {
			return err
		}
	}
	return nil

}

func (b *biliPlugin) sendRoomInfo(info *RoomInfo) error {
	db, err := b.env.GetDB()
	if err != nil {
		return err
	}
	record := &LiveRecord{}
	err = db.First(&record, info.Uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有记录，插入一条新记录并设置状态
			record = &LiveRecord{
				Uid:    int64(info.Uid),
				IsLive: IsLiving(info.LiveStatus),
			}
			err = db.Create(&record).Error
		}
		return err // 处理其他查询错误
	}
	if err != nil {
		return err
	}

	lastStatus := record.IsLive
	living := IsLiving(info.LiveStatus)
	change := lastStatus != living
	if change {
		// 状态有变化，更新数据库
		record.IsLive = living
		if err := db.Save(&record).Error; err != nil {
			return err
		}
	}
	if !change {
		return nil
	}
	if living {
		b.env.RangeBot(func(ctx *zero.Ctx) bool {
			msgChan := []message.MessageSegment{
				message.AtAll(),
				message.Text(fmt.Sprintf("\n@%s\n", info.Uname)),
				message.Image(info.Face),
				message.Text(fmt.Sprintf("\n%s\n%s\n", b.conf.randChoseLiveTips(), info.Title)),
				message.Image(info.CoverFromUser),
				message.Text(fmt.Sprintf("\nhttps://live.bilibili.com/%d", info.RoomId)),
			}
			for _, group := range b.conf.Groups {
				gopool.Go(func() {
					ctx.SendGroupMessage(group, msgChan)
				})
			}
			return true
		})
		return nil
	}
	if !living && b.conf.SendOff {
		b.env.RangeBot(func(ctx *zero.Ctx) bool {
			msgChan := []message.MessageSegment{
				message.Text(fmt.Sprintf("\n@%s\n", info.Uname)),
				message.Image(info.Face),
				message.Text(fmt.Sprintf("\n%s", b.conf.randChoseOffTips())),
			}
			for _, group := range b.conf.Groups {
				gopool.Go(func() {
					ctx.SendGroupMessage(group, msgChan)
				})
			}
			return true
		})
		return nil
	}

	return nil
}
