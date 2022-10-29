package config

import (
	"os"
)

var (
	APIPort     = SetEnv("APIPort", ":8080")
	APIKey      = SetEnv("APIKey", "AppSuberb-WAW")
	TokenSecret = "SuperbMIFTAH"
)

func SetEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}
