package store

import (
	"context"
	"jonghong/internal/pkg/model"

	"gorm.io/gorm"
)

type PaymentStore interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord *model.PaymentRecordM) error
	UpdatePaymentRecord(ctx context.Context, paymentRecord *model.PaymentRecordM) error
	GetPaymentRecord(ctx context.Context, paymentNo string, userId int) (*model.PaymentRecordM, error)
	GetPaymentRecordsByUser(ctx context.Context, userId int) ([]model.PaymentRecordM, error)
}

type payment struct {
	db *gorm.DB
}

// 接口实现

func (p *payment) CreatePaymentRecord(ctx context.Context, paymentRecord *model.PaymentRecordM) error {
	return p.db.Create(paymentRecord).Error
}

func (p *payment) UpdatePaymentRecord(ctx context.Context, paymentRecord *model.PaymentRecordM) error {

	return p.db.Model(&model.PaymentRecordM{}).Where("paymentNo = ?", paymentRecord.PaymentNo).Updates(paymentRecord).Error
}

func (p *payment) GetPaymentRecord(ctx context.Context, paymentNo string, userId int) (*model.PaymentRecordM, error) {

	var paymentRecord model.PaymentRecordM

	err := p.db.Where("paymentNo = ?", paymentNo).Where("userId = ?", userId).First(&paymentRecord).Error

	return &paymentRecord, err
}

func (p *payment) GetPaymentRecordsByUser(ctx context.Context, userId int) ([]model.PaymentRecordM, error) {

	var paymentRecords []model.PaymentRecordM
	err := p.db.Where("userId = ?", userId).Order("createTime desc").Find(&paymentRecords).Error

	return paymentRecords, err
}

// 验证

func newPayment(db *gorm.DB) PaymentStore {
	return &payment{db: db}
}
