package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	EnableHTTP  bool
}

func LoadConfig() Config {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	addr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if addr == "" {
		addr = ":8080"
	}
	enable := true
	if raw, ok := os.LookupEnv("HTTP_ENABLE"); ok {
		b, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			log.Printf("Flag error %v", err)
		} else {
			enable = b
		}
	} else if dsn == "" {
		enable = false
	}

	return Config{DatabaseURL: dsn, HTTPAddr: addr, EnableHTTP: enable}
}

func MaskDSN(dsn string) string {
	if dsn == "" {
		return ""
	}

	schemeIdx := strings.Index(dsn, "://")
	start := 0
	if schemeIdx >= 0 {
		start = schemeIdx + 3
	}

	atRel := strings.IndexByte(dsn[start:], '@')
	if atRel < 0 {
		return dsn
	}
	at := start + atRel

	userinfo := dsn[start:at]

	colon := strings.IndexByte(userinfo, ':')
	if colon < 0 {
		return dsn
	}

	user := userinfo[:colon]
	prefix := dsn[:start]
	suffix := dsn[at:]
	return prefix + user + ":****" + suffix
}
