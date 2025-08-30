package utils

import (
	"net/http"
	"strings"

	"github.com/akramboussanni/marchive/config"
	"github.com/akramboussanni/marchive/internal/model"
)

func cookieOp(name, value, path string, maxAge int) *http.Cookie {
	configuredDomain := config.App.Domain

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	}

	// Set domain for non-localhost domains
	// Note: localhost and 127.0.0.1 work without domain, but other local IPs might need it
	if configuredDomain != "localhost" && configuredDomain != "127.0.0.1" && !strings.HasPrefix(configuredDomain, "192.168.") && !strings.HasPrefix(configuredDomain, "10.") && !strings.HasPrefix(configuredDomain, "172.") {
		// Remove port if present for cookie domain
		domainWithoutPort := configuredDomain
		if colonIndex := strings.LastIndex(configuredDomain, ":"); colonIndex > 0 {
			domainWithoutPort = configuredDomain[:colonIndex]
		}
		cookie.Domain = domainWithoutPort
	}

	return cookie
}

func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, cookieOp("session", token, "/", int(config.App.JwtExpirations[string(model.CredentialJwt)])))
}

func SetRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, cookieOp("refresh", token, "/auth/refresh", int(config.App.JwtExpirations[string(model.RefreshJwt)])))
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, cookieOp("session", "", "/", -1))
}

func ClearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, cookieOp("refresh", "", "/auth/refresh", -1))
}

func ClearAllCookies(w http.ResponseWriter) {
	ClearSessionCookie(w)
	ClearRefreshCookie(w)
}
