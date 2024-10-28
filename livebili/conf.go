package livebili

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// 订阅的uid
	Uids []int64 `yaml:"uids" mapstructure:"uids"`
	// 检查直播间状态的间隔时间，单位为秒
	CheckLiveDuration int `yaml:"check_live_duration" mapstructure:"check_live_duration"`
	// 直播开始时推送的提示语 (精确到时分秒,用hh作为时，mm作为分，ss作为秒的占位符)
	LiveTips []string `yaml:"live_tips" mapstructure:"live_tips"`
	// 直播结束时推送的提示语 (精确到时分秒,用hh作为时，mm作为分，ss作为秒的占位符)
	OffTips []string `yaml:"off_tips" mapstructure:"off_tips"`
	// 下播是否推送
	SendOff bool `yaml:"send_off" mapstructure:"send_off"`
	// 检查动态的间隔时间，单位为秒
	CheckDynamicDuration int `yaml:"check_dynamic_duration" mapstructure:"check_dynamic_duration"`
	// 检查粉丝数的间隔时间，单位为秒(建议0.5-1小时查一次)
	CheckFollowerDuration int `yaml:"check_follower_duration" mapstructure:"check_follower_duration"`
	// b站的cookies
	Cookies string `yaml:"cookies" mapstructure:"cookies"`
}

func (c *Config) randChoseLiveTips() string {
	return c.LiveTips[rand.Intn(len(c.LiveTips))]
}
func (c *Config) randChoseOffTips() string {
	return c.OffTips[rand.Intn(len(c.OffTips))]
}

// 格式化时隔hh时mm分ss秒
func formatDuration(s string, dur time.Time) string {
	now := time.Now()
	sub := now.Sub(dur)
	// 将时间间隔转换为小时、分钟和秒
	hours := int(sub.Hours())
	minutes := int(sub.Minutes()) % 60
	seconds := int(sub.Seconds()) % 60
	s = strings.Replace(s, "hh", strconv.Itoa(hours), -1)
	s = strings.Replace(s, "mm", strconv.Itoa(minutes), -1)
	s = strings.Replace(s, "ss", strconv.Itoa(seconds), -1)
	return s
}
