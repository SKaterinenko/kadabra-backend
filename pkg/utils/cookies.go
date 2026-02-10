package utils

import (
	"kadabra/internal/core/config"
	"net/http"
)

// SetAuthCookies Установка токенов в cookies
func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string, cfg *config.Config) {
	// Access Token Cookie коротко живущий
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,                 // Защита от XSS
		Secure:   true,                 // Только HTTPS (в production)
		SameSite: http.SameSiteLaxMode, // Защита от CSRF
		MaxAge:   int(cfg.JWTAccessExpiration.Seconds()),
	})

	// Refresh Token Cookie долго живущий
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "api/auth/refresh", // Только для refresh endpoint
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(cfg.JWTRefreshExpiration.Seconds()),
	})
}

// ClearAuthCookies Очистка cookies при логауте
func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Удаляем cookie
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "api/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
