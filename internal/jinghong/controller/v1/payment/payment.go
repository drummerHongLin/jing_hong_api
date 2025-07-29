package payment

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/log"
	v1 "jonghong/pkg/api/jinghong/v1"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	b biz.IBiz
}

func NewPaymentController(ds store.IStore) PaymentController {
	return PaymentController{b: biz.NewBiz(ds)}
}

func (pc *PaymentController) CreateNewPaymentRecord(c *gin.Context) {
	log.C(c).Infow("Create new payment record function called")
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	var r v1.NewPaymentRecordRequest
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	err = pc.b.PaymentBiz().CreatePaymentRecord(c, &r, user.ID)

	core.WriteResponse(c, err, nil)

}

func (pc *PaymentController) UpdatePaymentRecord(c *gin.Context) {
	log.C(c).Infow("Update payment record function called")
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	var r v1.NewPaymentRecordRequest
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	err = pc.b.PaymentBiz().UpdatePaymentRecord(c, &r, user.ID)

	core.WriteResponse(c, err, nil)

}

func (pc *PaymentController) GetPaymentRecordByNo(c *gin.Context) {
	log.C(c).Infow("Get message by No function called")
	paymentNo := c.Param("paymentNo")

	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	record, err := pc.b.PaymentBiz().GetPaymentRecordByNo(c, paymentNo, user.ID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, record)

}

func (pc *PaymentController) GetPaymentRecordsById(c *gin.Context) {
	log.C(c).Infow("Get messages by Id function called")

	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	record, err := pc.b.PaymentBiz().GetPaymentRecordsById(c, user.ID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, record)

}
