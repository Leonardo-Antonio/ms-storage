package req

import "net/http"

func IsTLS(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
