package utils

import (
	"net/http"
	"time"
)

const DefaultTimeout = 10 * time.Second

func GetPort(r *http.Request) string {
	p := r.URL.Port()
	if p != "" {
		return p
	}

	switch r.URL.Scheme {
	case "http":
		return "80"
	default:
		return "443"
	}
}
