package auth

import (
	"github.com/akramboussanni/gocode/internal/jwt"
	"github.com/akramboussanni/gocode/internal/model"
)

func GenerateLogin(jwtToken jwt.Jwt) model.LoginTokens {
	sessionToken := jwtToken.WithType(model.CredentialJwt).GenerateToken()
	refreshToken := jwtToken.WithType(model.RefreshJwt).GenerateToken()

	return model.LoginTokens{
		Session: sessionToken,
		Refresh: refreshToken,
	}
}
