package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	// compose RegisteredClaims
	jwt.RegisteredClaims

	EmailVerified     bool                `json:"email_verified,omitempty"`
	Name              string              `json:"name,omitempty"`
	PreferredUsername string              `json:"preferred_username,omitempty"`
	GivenName         string              `json:"given_name,omitempty"`
	FamilyName        string              `json:"family_name,omitempty"`
	Email             string              `json:"email,omitempty"`
	ResourceAccess    interface{}         `json:"resource_access,omitempty"`
	RealmAccess       map[string][]string `json:"realm_access,omitempty"`
	Groups            []string            `json:"groups,omitempty"`

	// custom properties
	AuthType string
}

// NewRSAAuth returns a jwt.Keyfunc with empty key. Useful for checking token expiration.
func NewRSAAuth() (jwt.SigningMethod, jwt.Keyfunc) {
	return jwt.SigningMethodRS256, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	}
}

func NewHMACAuth(key []byte) (jwt.SigningMethod, jwt.Keyfunc) {
	return jwt.SigningMethodHS256, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}
}

// jwtCreds is JWT credentials for gRPC.
// implements gRPC credentials.PerRPCCredentials.
type jwtCreds struct {
	token    string
	authType string
}

func NewPerRPCInternalCredentials(token string) jwtCreds {
	return NewPerRPCCredentials("internalkey", token)
}

func NewPerRPCCredentials(authType, token string) jwtCreds {
	return jwtCreds{token: token, authType: authType}
}

func (j jwtCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.authType + " " + j.token,
	}, nil
}

func (j jwtCreds) RequireTransportSecurity() bool {
	return false
}
