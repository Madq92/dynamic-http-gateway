package utils

var isEnvLocal bool

// IsEnvLocal show whether executed locally
func IsEnvLocal() bool {
	return isEnvLocal
}

func init() {
	//isEnvLocal = os.Getenv("DHG_ENV") == "local"
	isEnvLocal = true
}
