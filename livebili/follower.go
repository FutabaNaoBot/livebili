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
	"time"
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
	r, err := b.checkFollower(uid)
	if err != nil {
		return err
	}

	db, err := b.env.GetDB()
	if err != nil {
		return err
	}
	now := time.Now()
	record := &FollowerRecord{Uid: uid}
	err = db.First(&record, uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有记录，插入一条新记录并设置状态
			record = &FollowerRecord{
				Uid:                uid,
				LastUpdate:         now,
				LastUpdateFollower: r.Data.Follower,
			}
			err = db.Create(&record).Error
		}
		return err // 处理其他查询错误
	}

	mode := b.followerNotifyMode(r.Data.Follower, record)
	if mode == NotNotify {
		return nil
	}
	// 保存变更
	err = FollowerRecord{
		Uid:                uid,
		LastUpdateFollower: r.Data.Follower,
		LastUpdate:         now,
	}.Save(db)
	if err != nil {
		return err
	}

	live, err := b.checkLive([]int64{uid})
	if err != nil {
		return err
	}
	var nickName string
	for _, roomInfo := range live.Data {
		nickName = roomInfo.Uname
		break
	}

	switch mode {
	case TimeReached:
		return b.onTimeReached(r.Data.Follower, record, nickName)
	case FollowerChange:
		return b.onFollowerChange(r.Data.Follower, record, nickName)
	case SpecialNumber:
		return b.onSpecialNumber(r.Data.Follower, record, nickName)
	case AroundSpecialNumber:
		return b.onAroundSpecialNumber(r.Data.Follower, record, nickName)
	default:
		return fmt.Errorf("unknown mode %d", mode)
	}

}

func (b *biliPlugin) checkFollower(uid int64) (r FollowerResp, err error) {
	resp, err := request.DoGet(fmt.Sprintf("https://api.bilibili.com/x/relation/stat?&vmid=%d", uid), b.conf.Cookies)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	r = FollowerResp{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return r, err
	}
	if r.Code != 0 {
		return r, fmt.Errorf("code: %d,msg: %s", r.Code, r.Message)
	}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type NotifyMode int

// 特殊数字步长(一万)
const specialNumStep = 10000

const (
	// NotNotify 无需通知
	NotNotify NotifyMode = iota
	// TimeReached 达到设定的通知时间
	TimeReached
	// FollowerChange 达到设定的粉丝变动数
	FollowerChange
	// SpecialNumber 达到特殊数字
	SpecialNumber
	// AroundSpecialNumber 在特殊数字周围
	AroundSpecialNumber
)

func (b *biliPlugin) followerNotifyMode(follower int, record *FollowerRecord) NotifyMode {
	if follower == record.LastUpdateFollower {
		return NotNotify
	}

	nowStep := follower / specialNumStep
	lastStep := record.LastUpdateFollower / specialNumStep
	if nowStep > lastStep {
		return SpecialNumber
	}
	if nowStep == lastStep {
		remain := specialNumStep - follower%specialNumStep
		// 与X万粉相差100以内
		if remain <= 100 {
			return AroundSpecialNumber
		}
	}

	now := time.Now()
	delta := follower - record.LastUpdateFollower

	if abs(delta)-b.conf.FollowerNotifyEach >= 0 {
		// 达到设定的每X个通知
		return FollowerChange
	} else if now.Sub(record.LastUpdate) > time.Duration(b.conf.FollowerNotifyDuration)*time.Second {
		// 超过设定的通知时间
		return TimeReached
	}

	return NotNotify
}

func (b *biliPlugin) onTimeReached(follower int, record *FollowerRecord, nickName string) error {
	return b.onFollowerChange(follower, record, nickName)
}

func (b *biliPlugin) onFollowerChange(follower int, record *FollowerRecord, nickName string) error {
	delta := follower - record.LastUpdateFollower
	var tips string
	if delta > 0 {
		tips = fmt.Sprintf("涨粉了！！🔺%d ", delta)
	} else {
		tips = fmt.Sprintf("掉粉了...🔻%d ", -delta)
	}
	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", nickName)),
		message.Text(tips),
		message.Text(fmt.Sprintf("%d → %d", record.LastUpdateFollower, follower)),
	)

	for ctx := range b.env.RangeBot {
		for gid := range b.groups.RangeGroup {
			ctx.SendGroupMessage(gid, msgChain)
		}
	}
	return nil
}

func (b *biliPlugin) onSpecialNumber(follower int, record *FollowerRecord, nickName string) error {
	step := follower / specialNumStep
	var tips string
	if step == 1 {
		tips = "🍾🎉万粉达成！！！🍾🎉"
	} else {
		tips = fmt.Sprintf("🍾🎉%d万粉达成！！！🍾🎉", step)
	}
	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", nickName)),
		message.Text(tips),
		message.Text(fmt.Sprintf("%d → %d", record.LastUpdateFollower, follower)),
	)
	for ctx := range b.env.RangeBot {
		for gid := range b.groups.RangeGroup {
			ctx.SendGroupMessage(gid, msgChain)
		}
	}
	return nil

}

func (b *biliPlugin) onAroundSpecialNumber(follower int, record *FollowerRecord, nickName string) error {
	nextStep := (follower / specialNumStep) + 1
	remain := (nextStep * specialNumStep) - follower
	var tips string
	if nextStep == 1 {
		tips = fmt.Sprintf("🎉距万粉剩余%d", remain)
	} else {
		tips = fmt.Sprintf("🎉距%d万粉剩余%d", nextStep, remain)
	}
	var msgChain chain.MessageChain
	msgChain.Split(
		message.Text(fmt.Sprintf("@%s", nickName)),
		message.Text(tips),
		message.Text(fmt.Sprintf("%d → %d", record.LastUpdateFollower, follower)),
	)
	for ctx := range b.env.RangeBot {
		for gid := range b.groups.RangeGroup {
			ctx.SendGroupMessage(gid, msgChain)
		}
	}
	return nil
}
