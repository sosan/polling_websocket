package repos

import (
	"polling_websocket/pkg/infra/tokenrepo"
	"time"
)

type AuthService interface {
	GenerateAccessToken() (*string, error)
	GetCachedActionUserAccessToken() *string
	VerifyActionUserToken(token string) (isOk bool, err error)
	VerifyUserToken(userToken string) (bool, bool)
}

type AuthRepository interface {
	GenerateIntrospectJWT(duration time.Duration) string
	GenerateAccessToken() (string, error)
	VerifyActionUserToken(token string) (bool, error)
	verifyWithIDProvider(token *tokenrepo.Token) (bool, error)
	VerifyUserToken(userToken string) (bool, bool)
}

type JWTGenerator interface {
	GenerateActionUserAssertionJWT(duration time.Duration) (string, error)
	GenerateAppInstrospectJWT(duration time.Duration) (string, error)
}

type ZitadelClient interface {
	GenerateActionUserAccessToken(jwt string) (*string, time.Duration, error)
	// GenerateActionUserAccessToken(jwt string) (*string, time.Duration, error)
	ValidateUserToken(userToken, introspectJWT string) (bool, int64, error)
	ValidateActionUserAccessToken(userToken, introspectJWT *string) (bool, error)
}

type TokenRepository interface {
	SaveToken(accessToken *string, expiresIn *time.Duration) error
	GetToken() (*tokenrepo.Token, error)
}
