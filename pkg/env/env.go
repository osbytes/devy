package env

import "os"

var Env string

const EnvLocal string = "local"
const EnvProd string = "prod"

func GetString(k string, d string) string {
	v, exists := os.LookupEnv(k)
	if !exists {
		return d
	}

	return v
}

func IsLocal() bool {
	return Env == EnvLocal
}

func IsProd() bool {
	return Env == EnvProd
}
