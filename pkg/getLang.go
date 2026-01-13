package pkg

import "net/http"

func GetLang(r *http.Request) string {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "ru"
	}
	return lang[:2]
}
