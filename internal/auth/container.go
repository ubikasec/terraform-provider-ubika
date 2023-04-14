package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

// ContainerAuthConfig define container authentication configuration.
type ContainerAuthConfig struct {
	Type Type   `json:"type"`
	Key  []byte `json:"key"`

	BaseURL string `json:"base_url"`

	AccessToken string `json:"access_token"`
}

func NewContainerAuthConfig(key, url string) *ContainerAuthConfig {
	ac := ContainerAuthConfig{
		Type: containerType,
		Key:  []byte(key),
	}
	ac.setURL(url)
	return &ac
}

// Login authenticates user with the authentication server.
func (ac *ContainerAuthConfig) Login() error {
	method, _ := NewHMACAuth(ac.Key)
	token := jwt.NewWithClaims(method, Claims{})
	ss, err := token.SignedString(ac.Key)
	if err == nil {
		ac.AccessToken = ss
	}
	return err
}

// Renew renews AccessToken.
func (ac *ContainerAuthConfig) Renew() error {
	return nil
}

// GetToken returns an access token for the current user.
func (ac *ContainerAuthConfig) GetToken() string {
	return ac.AccessToken
}

// Valid returns true if current AccessToken is valid.
func (ac *ContainerAuthConfig) Valid() bool {
	_, keyFunc := NewHMACAuth(ac.Key)
	_, err := jwt.Parse(ac.AccessToken, keyFunc)
	if err != nil {
		return false
	} else {
		return true
	}
}

// GetType returns the type of the Auth Configuration.
func (ac *ContainerAuthConfig) GetType() Type {
	return containerType
}

func (ac *ContainerAuthConfig) setURL(url string) {
	ac.BaseURL = url
}
