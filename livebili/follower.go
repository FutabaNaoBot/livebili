package livebili

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kohmebot/livebili/request"
	"github.com/kohmebot/pkg/chain"
	"github.com/kohmebot/pkg/gopool"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"io"
)

func (b *biliPlugin) doCheckFollower() error {
	var uids []int64
	for _, uid := range b.conf.Uids {
		uids = append(uids, uid)
	}
	errChan := make(chan error, len(uids))
	defer close(errChan)
	for _, uid := range uids {
		gopool.Go(func() {
			errChan <- b.doCheckOneFollower(uid)
		})
	}
	var err error
	for i := 0; i < len(uids); i++ {
		doErr := <-errChan
		if doErr != nil {
			err = errors.Join(err, fmt.Errorf("checkOneFollowerr %w", doErr))
		}
	}
	return err
}

func (b *biliPlugin) doCheckOneFollower(uid int64) error {
	resp, err := request.DoGet(fmt.Sprintf("https://api.bilibili.com/x/relation/stat?&vmid=%d", uid), b.conf.Cookies)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r := FollowerResp{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return fmt.Errorf("code: %d,msg: %s", r.Code, r.Message)
	}

	db, err := b.env.GetDB()
	if err != nil {
		return err
	}

	record := &FollowerRecord{Uid: uid}
	err = db.First(&record, uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有记录，插入一条新记录并设置状态
			record = &FollowerRecord{
				Uid:      uid,
				Follower: r.Data.Follower,
			}
			err = db.Create(&record).Error
		}
		return err // 处理其他查询错误
	}

	if record.Follower == r.Data.Follower {
		return nil
	}
	// 如果有粉丝数变更
	err = db.Save(&FollowerRecord{
		Uid:      uid,
		Follower: r.Data.Follower,
	}).Error
	if err != nil {
		return err
	}

	live, err := b.checkLive([]int64{uid})
	if err != nil {
		return err
	}
	var info RoomInfo
	for _, roomInfo := range live.Data {
		info = roomInfo
		break
	}
	var tips string
	delta := r.Data.Follower - record.Follower
	if delta > 0 {
		tips = fmt.Sprintf("涨粉了！！🔺%d ", delta)
	} else {
		tips = fmt.Sprintf("掉粉了...🔻%d ", delta)
	}

	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", info.Uname)),
		message.Text(tips),
		message.Text(fmt.Sprintf("%d → %d", record.Follower, r.Data.Follower)),
	)

	for ctx := range b.env.RangeBot {
		for gid := range b.groups.RangeGroup {
			ctx.SendGroupMessage(gid, msgChain)
		}
	}
	return nil
}
