package errno

// 定义常用的客户类型

var (
	OK                  = &Errno{HTTP: 200, Code: "", Message: ""}
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error."}
	ErrPageNotFound     = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page not found."}
	ErrUserAlreadyExist = &Errno{HTTP: 500, Code: "InternalError.UserAlreadyExist", Message: "User already exist."}
	ErrBind             = &Errno{HTTP: 400, Code: "InvalidParameter.BindError", Message: "Error occurred while binding the request body to the struct."}
	ErrInvalidParameter = &Errno{HTTP: 400, Code: "InvalidParameter", Message: "Parameter verification failed."}
	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &Errno{HTTP: 401, Code: "AuthFailure.SignTokenError", Message: "Error occurred while signing the JSON web token."}
	// ErrTokenInvalid 表示 JWT Token 格式错误.
	ErrTokenInvalid = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: "Token was invalid."}

	// ErrTokenInvalid 表示 JWT Token 超时.
	ErrTokenExpired = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: "Token was expired."}
	// ErrUnauthorized 表示请求没有被授权.
	ErrUnauthorized = &Errno{HTTP: 401, Code: "AuthFailure.Unauthorized", Message: "Unauthorized."}

	ErrUserNotFound      = &Errno{HTTP: 400, Code: "InternalError.UserNotFound", Message: "User not found."}
	ErrPasswordIncorrect = &Errno{HTTP: 400, Code: "InternalError.PasswordIncorrect", Message: "Password Incorrect."}

	ErrCodeNotExist = &Errno{HTTP: 401, Code: "AuthFailure.CodeNotExist", Message: "Code not exist."}
	ErrCodeExpired  = &Errno{HTTP: 401, Code: "AuthFailure.CodeInvalid", Message: "Code was expired."}
)
