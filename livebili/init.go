package livebili

import (
	"errors"
	"fmt"
	"github.com/kohmebot/plugin/pkg/chain"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"time"
)

func (b *biliPlugin) init() error {
	conf := Config{}
	db, err := b.env.GetDB()
	if err != nil {
		return err
	}

	err = b.env.GetConf(&conf)
	if err != nil {
		return err
	}
	b.conf = conf
	err = b.initData(db)
	if err != nil {
		return err
	}
	go b.tickerLive()
	go b.tickerDynamic()
	return nil
}

func (b *biliPlugin) initData(db *gorm.DB) error {
	err := db.AutoMigrate(&LiveRecord{})
	if err != nil {
		return err
	}
	for _, uid := range b.conf.Uids {
		record := &LiveRecord{Uid: uid}
		if err := db.Where("uid = ?", uid).First(record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				record.IsLive = false
				if err := db.Create(record).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

	}
	return db.AutoMigrate(&DynamicRecord{})
}

func (b *biliPlugin) tickerLive() {
	dur := time.Second * time.Duration(b.conf.CheckLiveDuration)
	t := time.NewTicker(dur)
	s := newErrorSender(b.sendError)
	for range t.C {
		logrus.Debugf("tock...%.2f sec", dur.Seconds())
		err := b.doCheckLive()
		if err != nil {
			logrus.Errorf("check live error: %v", err)
		}
		t.Reset(dur)
		s.Error(err)
	}

}

func (b *biliPlugin) tickerDynamic() {
	dur := time.Second * time.Duration(b.conf.CheckDynamicDuration)
	t := time.NewTicker(dur)
	s := newErrorSender(b.sendError)
	for range t.C {
		logrus.Debugf("tock...%.2f sec", dur.Seconds())
		err := b.doCheckDynamic()
		if err != nil {
			logrus.Errorf("check dynamic error: %v", err)
		}
		t.Reset(dur)
		s.Error(err)
	}
}

func (b *biliPlugin) sendError(err error) {
	b.env.RangeBot(func(ctx *zero.Ctx) bool {
		b.groups.RangeGroup(func(group int64) bool {
			var msgChain chain.MessageChain
			msgChain.Split(
				message.Text("我出错了喵！快帮我联系管理员喵！！"),
				message.Text(err.Error()),
				message.Text(fmt.Sprintf("from: %s", b.Name())),
			)
			ctx.SendGroupMessage(group, msgChain)
			return true
		})
		return true
	})
}
