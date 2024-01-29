package httputils

import (
	"dalkak/pkg/utils/securityutils"
	"net/http"
	"time"
)

func SetRefreshToken(w http.ResponseWriter, mode string, domain string, token string, tokenTime int64) {
	tokenExpires := time.Unix(tokenTime, 0)
	refreshTokenDuration := time.Duration(securityutils.RefreshTokenTTL) * time.Second
	expires := tokenExpires.Add(refreshTokenDuration)

	parsedDomain, err := ParseDomain(domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  isSecure := mode != "LOCAL"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    token,
		Expires:  expires,
		MaxAge:   securityutils.RefreshTokenTTL,
		SameSite: http.SameSiteStrictMode,
		Domain:   parsedDomain,
		HttpOnly: true,
		Secure:   isSecure,
	})
}
