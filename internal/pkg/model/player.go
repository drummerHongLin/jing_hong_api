package model

import "time"

type PlayConfigM struct {
	ID            int       `gorm:"column:id;" json:"id"`
	UserId        int       `gorm:"column:userId; not null" json:"userId"`
	User          UserM     `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlayConfigNo  string    `gorm:"column:playConfigNo;not null;" json:"playConfigNo"`
	Bpm           int       `gorm:"column:bpm;not null;" json:"bpm"`
	BeatNum       int       `gorm:"column:beatNum;not null;" json:"beatNum"`
	BeatNote      int       `gorm:"column:beatNote;not null;" json:"beatNote"`
	ReferenceBeat int       `gorm:"column:referenceBeat;not null;" json:"referenceBeat"`
	SubBeats      string    `gorm:"column:subBeats;not null;" json:"subBeats"`
	ConfigTitle   string    `gorm:"column:configTitle;not null;" json:"configTitle"`
	CreatedAt     time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (p *PlayConfigM) TableName() string {
	return "play_config"
}
