package v1

// 2. 登录返回数据格式 token
type EmailVerifiedResponse struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}

type EmailVerifingRequest struct {
	Code string `json:"code"  validate:"required,max=6,min=6"`
}
