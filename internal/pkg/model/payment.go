package model

import "time"

type PaymentRecordM struct {
	ID            int     `gorm:"column:id;" json:"id"`
	UserId        int     `gorm:"column:userId; not null" json:"userId"`
	User          UserM   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PaymentNo     string  `gorm:"column:paymentNo;not null;" json:"paymentNo"`
	Amount        float64 `gorm:"column:amount;not null;" json:"amount"`
	CreateTime    string  `gorm:"column:createTime;not null" json:"createTime"`
	PayTime       string  `gorm:"column:payTime;not null" json:"payTime"`
	PaymentStatus int     `gorm:"column:paymentStatus;not null" json:"paymentStatus"`
	Label         string  `gorm:"column:label;not null" json:"label"`
	// CreatedAt UpdatedAt 是保留关键字
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (u *PaymentRecordM) TableName() string {
	return "payment_record"
}

/*
因为前后端都是都是自己写的，所以先不用枚举

type status int

const (
	_ status = iota
	created
	canceled
	completed
	refunded
)

func (s status) Valid() bool {
	return s >= created && s <= refunded
}

func (s status) String() string {
	return [...]string{"_", "已创建", "已取消", "已完成", "已退款"}[s]
}

func ParseStatus(i int) (status, error) {
	s := status(i)
	if !s.Valid() {
		return -1, errno.InternalServerError.SetMessage("支付状态无效!")
	}
	return s, nil
}
*/
