package known

var (
	XRequestIDKey = "X-Request-ID"
	XUsernameKey  = "X-Username"
	XUserIdKey    = "X-UserId"
	HomeDir       string
)

func SetHomeDir(path string) {
	HomeDir = path
}
