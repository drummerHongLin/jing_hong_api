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
	UpdatePaymentRecord(ctx context.Context, r *v1.NewPaymentRecordRequest, userId int) error
	GetPaymentRecordByNo(ctx context.Context, paymentNo string, userId int) (*v1.GetPaymentRecordResponse, error)
	GetPaymentRecordsById(ctx context.Context, userId int) (*v1.GetPaymentRecordsResponse, error)
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
	error := p.ps.CreatePaymentRecord(ctx, &newPaymentRecord)

	return error
}

func (p *payment) UpdatePaymentRecord(ctx context.Context, r *v1.NewPaymentRecordRequest, userId int) error {

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

func (p *payment) GetPaymentRecordsById(ctx context.Context, userId int) (*v1.GetPaymentRecordsResponse, error) {
	paymentRecords, err := p.ps.GetPaymentRecordsByUser(ctx, userId)

	if err != nil {
		return nil, err
	}

	var recordsResponse []v1.PaymentRecord
	_ = copier.Copy(&recordsResponse, paymentRecords)

	return &v1.GetPaymentRecordsResponse{
		PaymentRecords: recordsResponse,
	}, nil

}
