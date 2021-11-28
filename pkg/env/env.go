package env

import "os"

func GetString(k string, d string) string {
	v, exists := os.LookupEnv(k)
	if !exists {
		return d
	}

	return v
}
