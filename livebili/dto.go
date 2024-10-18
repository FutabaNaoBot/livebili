package livebili

type LiveRecord struct {
	Uid    int64 `gorm:"primaryKey"`
	IsLive bool
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
