package httputils

import (
	"dalkak/config"
	"dalkak/pkg/dtos"
	"net/http"
	"time"
)

func SetCookieRefresh(w http.ResponseWriter, mode string, refreshToken string, tokenTime int64, domain string) error {
	tokenExpires := time.Unix(tokenTime, 0)
	refreshTokenDuration := time.Duration(config.RefreshTokenTTL) * time.Second
	refreshTokenexpires := tokenExpires.Add(refreshTokenDuration)

	parsedDomain, err := ParseDomain(domain)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "failed to parse domain",
		}
	}

	isSecure := mode != "LOCAL"
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    refreshToken,
		Expires:  refreshTokenexpires,
		MaxAge:   config.RefreshTokenTTL,
		SameSite: http.SameSiteLaxMode,
		Domain:   parsedDomain,
		HttpOnly: true,
		Secure:   isSecure,
	})
	return nil
}

func GetCookieRefresh(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "refresh token not found",
		}
	}
	return cookie.Value, nil
}

func DeleteCookieRefresh(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Path:   "/",
		MaxAge: -1,
	})
}
