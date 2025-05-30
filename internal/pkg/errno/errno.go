package errno

import (
	"fmt"
)

// 统一错误的返回信息，http错误码，内部错误码，错误信息

type Errno struct {
	HTTP    int
	Code    string
	Message string
}

// 对自定义的错误结构继承error

func (en *Errno) Error() string {
	return en.Message
}

// 设置消息方法

func (en *Errno) SetMessage(format string, args ...any) *Errno {
	en.Message = fmt.Sprintf(format, args...)
	return en
}

// 验证是否接口实现

var _ error = &Errno{}

// 从错误中解析出错误码和信息

func DecodeError(err error) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message
	}
	// 判断错误的类型是否是定义的结构体
	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
		return InternalServerError.HTTP, InternalServerError.Code, err.Error()
	}

}
