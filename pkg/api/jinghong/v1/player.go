package v1

import "time"

// create和update都用这个
type PlayConfig struct {
	PlayConfigNo  string `json:"playConfigNo"`
	Bpm           int    `json:"bpm"`
	BeatNum       int    `json:"beatNum"`
	BeatNote      int    `json:"beatNote"`
	ReferenceBeat int    `json:"referenceBeat"`
	SubBeats      string `json:"subBeats"`
	ConfigTitle   string `json:"configTitle"`
}

// 返回信息多携带创建和更新时间
type PlayConfigWithTime struct {
	PlayConfig
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
}

type GetPlayConfigsResponse struct {
	HasMore     bool                 `json:"hasMore"`
	PlayConfigs []PlayConfigWithTime `json:"playConfigs"`
}
