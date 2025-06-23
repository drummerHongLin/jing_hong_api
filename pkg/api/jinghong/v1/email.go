package v1

// 2. 登录返回数据格式 token
type EmailVerifiedResponse struct {
	Token string `json:"token"`
}

type EmailVerifingRequest struct {
	Code string `json:"code"  validate:"required,max=6,min=6"`
}
