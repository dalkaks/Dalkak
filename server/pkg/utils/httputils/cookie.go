package httputils

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/securityutils"
	"net/http"
	"time"
)

func SetCookieTokenPair(w http.ResponseWriter, mode string, authTokens *dtos.AuthTokens, tokenTime int64, domain string) {
	tokenExpires := time.Unix(tokenTime, 0)

	accessTokenDuration := time.Duration(securityutils.AccessTokenTTL) * time.Second
	accessTokenExpires := tokenExpires.Add(accessTokenDuration)

	refreshTokenDuration := time.Duration(securityutils.RefreshTokenTTL) * time.Second
	refreshTokenexpires := tokenExpires.Add(refreshTokenDuration)

	parsedDomain, err := ParseDomain(domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isSecure := mode != "LOCAL"

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    authTokens.AccessToken,
		Expires:  accessTokenExpires,
		MaxAge:   securityutils.AccessTokenTTL,
		SameSite: http.SameSiteStrictMode,
		Domain:   parsedDomain,
		HttpOnly: true,
		Secure:   isSecure,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    authTokens.RefreshToken,
		Expires:  refreshTokenexpires,
		MaxAge:   securityutils.RefreshTokenTTL,
		SameSite: http.SameSiteStrictMode,
		Domain:   parsedDomain,
		HttpOnly: true,
		Secure:   isSecure,
	})
}
