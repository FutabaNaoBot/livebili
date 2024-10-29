package livebili

import (
	"gorm.io/gorm"
	"sort"
	"time"
)

type FollowerRecord struct {
	Uid int64 `gorm:"primaryKey"`
	// 上次更新时间
	LastUpdate time.Time
	// 上次更新的粉丝数
	LastUpdateFollower int
}

func (f FollowerRecord) Save(db *gorm.DB) error {
	return db.Save(f).Error
}

type LiveRecord struct {
	Uid    int64 `gorm:"primaryKey"`
	IsLive bool
	// 上次直播时间
	LastLiveTime time.Time
	// 上次下播时间
	LastOffTime time.Time
}

type LiveResp struct {
	Code    int                 `json:"code"`
	Msg     string              `json:"msg"`
	Message string              `json:"message"`
	Data    map[string]RoomInfo `json:"data"`
}

type RoomInfo struct {
	Title         string `json:"title"`
	RoomId        int    `json:"room_id"`
	Uname         string `json:"uname"`
	Face          string `json:"face"`
	CoverFromUser string `json:"cover_from_user"`
	LiveStatus    int    `json:"live_status"`
	Uid           int    `json:"uid"`
}

func IsLiving(status int) bool {
	return status == 1
}

type DynamicRecord struct {
	Uid         int64 `gorm:"primaryKey"`
	LastPubTime int64
}

type DynamicResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    DynamicData `json:"data"`
}

func (d *DynamicResp) TopReSort() {
	topIdx := -1
	for idx, item := range d.Data.Items {
		if item.Modules.ModuleTag.Text == "置顶" {
			topIdx = idx
		}
		break
	}
	if topIdx < 0 {
		return
	}
	sort.Slice(d.Data.Items, func(i, j int) bool {
		return d.Data.Items[i].Modules.PutTs > d.Data.Items[j].Modules.PutTs
	})

}

type DynamicData struct {
	HasMore bool      `json:"has_more"`
	Items   []Dynamic `json:"items"`
}

type Dynamic struct {
	IdStr   string         `json:"id_str"`
	Type    string         `json:"type"`
	Modules DynamicModules `json:"modules"`
}

type DynamicModules struct {
	ModuleAuthor  `json:"module_author"`
	ModuleDynamic `json:"module_dynamic"`
	ModuleTag     `json:"module_tag"`
}

// ModuleTag 置顶信息
type ModuleTag struct {
	// 置顶动态出现这个对象，否则没有
	Text string `json:"text"`
}

type ModuleAuthor struct {
	// UP主名称
	// 剧集名称
	// 合集名称
	Name string `json:"name"`
	// x分钟前
	// x小时前
	// 昨天
	PubTime string `json:"pub_time"`
	// Unix秒级时间戳
	PutTs int64 `json:"pub_ts"`
	// 投稿了视频
	// 直播了
	// 投稿了文章
	// 更新了合集
	// 与他人联合创作
	// 发布了动态视频
	// 投稿了直播回放
	PubAction string `json:"pub_action"`
}

type ModuleDynamic struct {
	Major `json:"major"`
	Desc  `json:"desc"`
}

type Desc struct {
	// 动态的文字内容
	Text string `json:"text"`
}

// Major 动态主体对象
type Major struct {
	// 动态主体类型
	Type    string       `json:"type"`
	Archive MajorArchive `json:"archive"`
	Draw    MajorDraw    `json:"draw"`
}

// MajorArchive 视频信息
type MajorArchive struct {
	Cover   string `json:"cover"`
	Title   string `json:"title"`
	JumpUrl string `json:"jump_url"`
}

// MajorDraw 带图动态
type MajorDraw struct {
	// 图片信息列表
	Items []DrawItem `json:"items"`
}

type DrawItem struct {
	// 图片URL
	Src string `json:"src"`
}

type FollowerResp struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    RelationData `json:"data"`
}

type RelationData struct {
	// uid
	Mid int `json:"mid"`
	// 关注数
	Following int `json:"following"`
	// 粉丝数
	Follower int `json:"follower"`
}
