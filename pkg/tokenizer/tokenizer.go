package tokenizer

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type JWTConfig struct {
	Secret string
	Expr   time.Duration
}

type MicroserviceClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type BasicAuthenticationConfig struct {
	Username string
	Password string
}

type ITokenizer interface {
	VerifyJWTToken(token string) (*jwt.Token, error)
	IsBasicAuthorized(token string) bool
}

type tokenizer struct {
	jwtConfig       JWTConfig
	basicAuthConfig BasicAuthenticationConfig
}

func NewTokenizer(jc JWTConfig, bc BasicAuthenticationConfig) ITokenizer {
	return &tokenizer{
		jwtConfig:       jc,
		basicAuthConfig: bc,
	}
}

func (t *tokenizer) VerifyJWTToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.jwtConfig.Secret), nil
	})
}

func (t *tokenizer) IsBasicAuthorized(token string) bool {
	decode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	return string(decode) == fmt.Sprintf("%s:%s", t.basicAuthConfig.Username, t.basicAuthConfig.Password)
}
