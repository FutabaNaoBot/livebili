package livebili

import "math/rand"

type Config struct {
	// 订阅的uid
	Uids []int64 `yaml:"uids" mapstructure:"uids"`
	// 检查直播间状态的间隔时间，单位为秒
	CheckLiveDuration int `yaml:"check_live_duration" mapstructure:"check_live_duration"`
	// 直播开始时推送的提示语
	LiveTips []string `yaml:"live_tips" mapstructure:"live_tips"`
	// 直播结束时推送的提示语
	OffTips []string `yaml:"off_tips" mapstructure:"off_tips"`
	// 下播是否推送
	SendOff bool `yaml:"send_off" mapstructure:"send_off"`
	// 检查动态的间隔时间，单位为秒
	CheckDynamicDuration int `yaml:"check_dynamic_duration" mapstructure:"check_dynamic_duration"`
}

func (c *Config) randChoseLiveTips() string {
	return c.LiveTips[rand.Intn(len(c.LiveTips))]
}
func (c *Config) randChoseOffTips() string {
	return c.OffTips[rand.Intn(len(c.OffTips))]
}
