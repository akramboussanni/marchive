package model

type JwtType string

const (
	CredentialJwt JwtType = "credential"
	RefreshJwt    JwtType = "refresh"
)
