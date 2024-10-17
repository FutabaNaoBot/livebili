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
