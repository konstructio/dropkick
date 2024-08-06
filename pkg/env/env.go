package env

import (
	"os"
	"strings"
)

func GetFirstNotEmpty(keys ...string) string {
	for _, v := range keys {
		if s := strings.TrimSpace(os.Getenv(v)); s != "" {
			return s
		}
	}

	return ""
}
