package httputils

import (
	"dalkak/config"
	"net/http"
	"time"
)

func SetCookieRefresh(w http.ResponseWriter, mode string, refreshToken string, tokenTime int64, domain string) {
	tokenExpires := time.Unix(tokenTime, 0)
	refreshTokenDuration := time.Duration(config.RefreshTokenTTL) * time.Second
	refreshTokenexpires := tokenExpires.Add(refreshTokenDuration)

	parsedDomain, err := ParseDomain(domain)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
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
}

func GetCookieRefresh(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", &AppError{http.StatusUnauthorized, "invalid refresh token"}
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
