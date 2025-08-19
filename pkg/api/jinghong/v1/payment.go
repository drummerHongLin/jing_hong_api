package v1

type NewPaymentRecordRequest struct {
	PaymentNo    string  `json:"paymentNo"`
	Amount       float64 `json:"amount"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	CreateTime   string  `json:"createTime"`
	ProductName  string  `json:"productName"`
	AccountToken string  `json:"accountToken"`
}

type UpdatePaymentRecordRequest struct {
	PaymentNo     string `json:"paymentNo"`
	PayTime       string `json:"payTime"`
	PayStatus     int    `json:"payStatus"`
	TransactionId string `json:"transactionId"`
}

type PaymentRecord struct {
	PaymentNo     string  `json:"paymentNo"`
	Amount        float64 `json:"amount"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	CreateTime    string  `json:"createTime"`
	PayTime       string  `json:"payTime"`
	PayStatus     int     `json:"payStatus"`
	ProductName   string  `json:"productName"`
	AccountToken  string  `json:"accountToken"`
	TransactionId string  `json:"transactionId"`
}

type GetPaymentRecordResponse struct {
	PaymentRecord
}

type GetPaymentRecordsResponse struct {
	HasMore        bool            `json:"hasMore"`
	PaymentRecords []PaymentRecord `json:"paymentRecords"`
}
