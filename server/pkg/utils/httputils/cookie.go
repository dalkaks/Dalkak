package httputils

import (
	"dalkak/pkg/utils/securityutils"
	"net/http"
	"time"
)

func SetCookieRefresh(w http.ResponseWriter, mode string, refreshToken string, tokenTime int64, domain string) {
	tokenExpires := time.Unix(tokenTime, 0)
	refreshTokenDuration := time.Duration(securityutils.RefreshTokenTTL) * time.Second
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
		MaxAge:   securityutils.RefreshTokenTTL,
		SameSite: http.SameSiteLaxMode,
		Domain:   parsedDomain,
		HttpOnly: true,
		Secure:   isSecure,
	})
}

func GetCookieRefresh(r *http.Request) string {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}
