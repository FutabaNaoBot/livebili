package livebili

import (
	"errors"
	"github.com/sirupsen/logrus"
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
	for range t.C {
		logrus.Debugf("tock...%.2f sec", dur.Seconds())
		err := b.doCheckLive()
		if err != nil {
			logrus.Errorf("check live error: %v", err)
		}
		t.Reset(dur)
	}

}

func (b *biliPlugin) tickerDynamic() {
	dur := time.Second * time.Duration(b.conf.CheckDynamicDuration)
	t := time.NewTicker(dur)
	for range t.C {
		logrus.Debugf("tock...%.2f sec", dur.Seconds())
		err := b.doCheckDynamic()
		if err != nil {
			logrus.Errorf("check dynamic error: %v", err)
		}
		t.Reset(dur)
	}
}
