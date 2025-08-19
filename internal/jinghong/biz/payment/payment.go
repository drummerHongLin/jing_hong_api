package payment

import (
	"context"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/model"
	v1 "jonghong/pkg/api/jinghong/v1"

	"github.com/jinzhu/copier"
)

type PaymentBiz interface {
	CreatePaymentRecord(ctx context.Context, r *v1.NewPaymentRecordRequest, userId int) error
	UpdatePaymentRecord(ctx context.Context, r *v1.UpdatePaymentRecordRequest, userId int) error
	GetPaymentRecordByNo(ctx context.Context, paymentNo string, userId int) (*v1.GetPaymentRecordResponse, error)
	GetPaymentRecordsById(ctx context.Context, userId int, offset int, limit int) (*v1.GetPaymentRecordsResponse, error)
	InsertPaymentRecord(ctx context.Context, r []v1.PaymentRecord, userId int) error
}

type payment struct {
	ps store.PaymentStore
}

func NewPaymentBiz(ps store.PaymentStore) PaymentBiz {
	return &payment{ps: ps}
}

func (p *payment) CreatePaymentRecord(ctx context.Context, r *v1.NewPaymentRecordRequest, userId int) error {

	var newPaymentRecord model.PaymentRecordM
	_ = copier.Copy(&newPaymentRecord, r)
	newPaymentRecord.UserId = userId
	newPaymentRecord.PayStatus = 0
	error := p.ps.CreatePaymentRecord(ctx, &newPaymentRecord)

	return error
}

func (p *payment) UpdatePaymentRecord(ctx context.Context, r *v1.UpdatePaymentRecordRequest, userId int) error {

	var newPaymentRecord model.PaymentRecordM
	_ = copier.Copy(&newPaymentRecord, r)
	newPaymentRecord.UserId = userId
	error := p.ps.UpdatePaymentRecord(ctx, &newPaymentRecord)

	return error
}

func (p *payment) GetPaymentRecordByNo(ctx context.Context, paymentNo string, userId int) (*v1.GetPaymentRecordResponse, error) {
	paymentRecord, err := p.ps.GetPaymentRecord(ctx, paymentNo, userId)

	if err != nil {
		return nil, err
	}

	var recordResponse *v1.GetPaymentRecordResponse

	_ = copier.Copy(recordResponse, paymentRecord)

	return recordResponse, nil

}

func (p *payment) GetPaymentRecordsById(ctx context.Context, userId int, offset int, limit int) (*v1.GetPaymentRecordsResponse, error) {
	paymentRecords, err := p.ps.GetPaymentRecordsByUser(ctx, userId, offset, limit)

	if err != nil {
		return nil, err
	}

	more, err := p.ps.GetPaymentRecordsByUser(ctx, userId, offset+limit, 1)

	if err != nil {
		return nil, err
	}

	hasMore := len(more) > 0

	var recordsResponse []v1.PaymentRecord
	_ = copier.Copy(&recordsResponse, paymentRecords)

	return &v1.GetPaymentRecordsResponse{
		HasMore:        hasMore,
		PaymentRecords: recordsResponse,
	}, nil

}

func (p *payment) InsertPaymentRecord(ctx context.Context, r []v1.PaymentRecord, userId int) error {

	var paymentRecords []model.PaymentRecordM
	_ = copier.Copy(&paymentRecords, &r)

	// 循环给每个记录添加userID
	for i := range paymentRecords {
		paymentRecords[i].UserId = userId
	}

	err := p.ps.InsertPaymentRecord(ctx, paymentRecords)

	return err

}
