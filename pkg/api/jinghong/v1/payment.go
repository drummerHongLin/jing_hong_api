package v1

type NewPaymentRecordRequest struct {
	PaymentRecord
}

type PaymentRecord struct {
	PaymentNo     string  `json:"paymentNo"`
	Amount        float64 `json:"amount"`
	CreateTime    string  `json:"createTime"`
	PayTime       string  `json:"payTime"`
	PaymentStatus int     `json:"paymentStatus"`
	Label         string  `json:"label"`
}

type GetPaymentRecordResponse struct {
	PaymentRecord
}

type GetPaymentRecordsResponse struct {
	PaymentRecords []PaymentRecord `json:"paymentRecords"`
}
