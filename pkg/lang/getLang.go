package lang

import (
	"net/http"
)

func GetLang(r *http.Request) string {
	lang := r.Header.Get("Accept-Language")
	if len(lang) <= 1 {
		lang = "ru"
	}
	return lang[:2]
}
