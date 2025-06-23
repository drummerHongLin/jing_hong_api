package v1

// 定义接口的数据规范

// 1. 登录传入数据格式 采用邮箱+密码登录
type LoginRequest struct {
	Username string `json:"username" validate:"alphanum,required,max=255,min=1"`
	Password string `json:"password" validate:"required,max=18,min=6"`
}

// 2. 登录返回数据格式 token
type LoginResponse struct {
	Token string `json:"token"`
}

// 3. 创建用户格式，头像和认证信息不是前端传入，后续在服务端修改
type CreateUserRequest struct {
	Username string `json:"username" validate:"alphanum,required,max=255,min=1"`
	Password string `json:"password" validate:"required,max=18,min=6"`
	Nickname string `json:"nickname" validate:"required,max=255,min=1"`
	Email    string `json:"email" validate:"required,email"`
}

// 4. 更改密码格式，更改密码操作需要验证token，所以不需要传入用户名
type ChangePasswordRequest struct {
	// 旧密码.
	OldPassword string `json:"oldPassword" validate:"required,max=18,min=6"`

	// 新密码.
	NewPassword string `json:"newPassword" validate:"required,max=18,min=6"`
}

// 5. 设置头像
// 头像上传至阿里云的oss服务
// 1）前端先向本服务端请求oss的签名信息
// 2）前端整合签名信息后将文件传入阿里云的oss
// 3）本服务端监听oss的回调
