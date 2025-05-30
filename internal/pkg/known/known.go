package known

var (
	XRequestIDKey = "X-Request-ID"
	XUsernameKey  = "X-Username"
	HomeDir       string
)

func SetHomeDir(path string) {
	HomeDir = path
}
