package errno

// 定义常用的客户类型

var (
	OK                  = &Errno{HTTP: 200, Code: "", Message: ""}
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "未知服务器错误！"}
	ErrPageNotFound     = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "该页面不存在！"}
	ErrUserAlreadyExist = &Errno{HTTP: 500, Code: "InternalError.UserAlreadyExist", Message: "用户名已存在！"}
	ErrBind             = &Errno{HTTP: 400, Code: "InvalidParameter.BindError", Message: "请求体错误！"}
	ErrInvalidParameter = &Errno{HTTP: 400, Code: "InvalidParameter", Message: "请求参数错误！"}
	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &Errno{HTTP: 401, Code: "AuthFailure.SignTokenError", Message: "签发token失败！"}
	// ErrTokenInvalid 表示 JWT Token 格式错误.
	ErrTokenInvalid = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: "无效token！"}

	// ErrTokenInvalid 表示 JWT Token 超时.
	ErrTokenExpired = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: "token已过期！"}
	// ErrUnauthorized 表示请求没有被授权.
	ErrUnauthorized = &Errno{HTTP: 401, Code: "AuthFailure.Unauthorized", Message: "用户未被授权！"}

	ErrUserNotFound = &Errno{HTTP: 401, Code: "InternalError.UserNotFound", Message: "用户不存在！"}

	ErrUserNotVerified = &Errno{HTTP: 401, Code: "InternalError.UserNotVerified", Message: "用户邮箱未被验证！"}

	ErrPasswordIncorrect = &Errno{HTTP: 401, Code: "InternalError.PasswordIncorrect", Message: "密码错误！"}

	ErrCodeNotExist = &Errno{HTTP: 401, Code: "AuthFailure.CodeNotExist", Message: "无效验证码！"}
	ErrCodeExpired  = &Errno{HTTP: 401, Code: "AuthFailure.CodeInvalid", Message: "验证码已过期！"}

	// 文件相关错误
	ErrNoFileUpdated     = &Errno{HTTP: 501, Code: "FileFailuer.NoFileUpdated", Message: "未找到相关文件！"}
	ErrFileTypeUnsupport = &Errno{HTTP: 501, Code: "FileFailuer.FileTypeUnsupport", Message: "不支持文件类型！"}

	FileSysError = &Errno{HTTP: 501, Code: "FileFailuer.FileSysError", Message: "未知文件系统异常！"}

	// 对话存储相关
	ErrSessionAlreadyExist = &Errno{HTTP: 500, Code: "InternalError.SessionAlreadyExist", Message: "会话ID已存在！"}

	// 枚举值相关

	ErrEnumOutOfRange = &Errno{HTTP: 500, Code: "EnumError.IndexOutOfRange", Message: "枚举值越界不存在"}
)
