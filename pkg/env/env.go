package env

import "os"

var Env string

const EnvLocal string = "local"
const EnvProd string = "prod"

func GetString(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

func IsLocal() bool {
	return Env == EnvLocal
}

func IsProd() bool {
	return Env == EnvProd
}
