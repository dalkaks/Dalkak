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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isSecure := mode != "LOCAL"
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    refreshToken,
		Expires:  refreshTokenexpires,
		MaxAge:   securityutils.RefreshTokenTTL,
		SameSite: http.SameSiteStrictMode,
		Domain:   parsedDomain,
		HttpOnly: true,
		Secure:   isSecure,
	})
}
