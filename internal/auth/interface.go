package auth

type Type string

type Authentifier interface {
	Login() error
	Renew() error
	Valid() bool
	GetToken() string
	GetType() Type
}

// BaseAuthConfig describes a base AuthConfig used to find the good implementation.
type BaseAuthConfig struct {
	Type Type `json:"type"`
}

const (
	keycloakType  Type = "keycloak"
	containerType Type = "container"
)

func (ac *BaseAuthConfig) Login() error     { return nil }
func (ac *BaseAuthConfig) Renew() error     { return nil }
func (ac *BaseAuthConfig) Valid() bool      { return false }
func (ac *BaseAuthConfig) GetToken() string { return "" }
func (ac *BaseAuthConfig) GetType() Type    { return Type("") }
