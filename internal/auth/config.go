package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var errAuthBaseFileNameEmpty = errors.New("auth base file name must not be empty")

type Config struct {
	CurrentContext string             `json:"current_context"`
	Contexts       map[string]Context `json:"contexts"`

	// path is the path to the configuration file
	path string `json:"-"`
}

type Context struct {
	AuthType   Type            `json:"auth_type"`
	AuthConfig json.RawMessage `json:"auth_config"`
}

// newConfig returns an empty Config.
func newConfig() Config {
	return Config{Contexts: make(map[string]Context)}
}

func Load(baseFileName string) (Config, error) {
	configPath := os.Getenv("APPSECCTL_CACHE_PATH")

	if configPath == "" {
		if baseFileName == "" {
			return Config{}, errAuthBaseFileNameEmpty
		}

		// search for usercache directory
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return Config{}, err
		}

		configPath = filepath.Join(cacheDir, baseFileName)
	}

	config := newConfig()
	config.path = configPath
	err := config.load()
	return config, err
}

func (c *Config) load() error {
	f, err := os.OpenFile(c.path, os.O_RDONLY, 0o600)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open config file: %s", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read config file: %s", err)
	}
	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}

// UseContext set current context to context name.
// Creates a new empty context if does not exist.
func (c *Config) UseContext(contextName string) {
	if _, ok := c.Contexts[contextName]; !ok {
		c.Contexts[contextName] = Context{}
	}
	c.CurrentContext = contextName
}

// Save save the auth configs to the disk.
func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// write auth cache file
	f, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open config file: %s", err)
	}

	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
		return err
	}

	return nil
}

// UpdateContext update the current context with a new Autentifier.
func (c *Config) UpdateContext(newAuth Authentifier) error {
	b, err := json.Marshal(newAuth)
	if err != nil {
		return err
	}

	if c.CurrentContext == "" {
		c.CurrentContext = "default"
	}

	c.Contexts[c.CurrentContext] = Context{
		AuthType:   newAuth.GetType(),
		AuthConfig: json.RawMessage(b),
	}

	return nil
}

func (c *Config) IsContext(contextName string) bool {
	_, ok := c.Contexts[contextName]
	return ok
}

// GetAuthConfig returns the auth configuration corresponding to the current context.
func (c *Config) GetAuthConfig(httpClient *http.Client) (Authentifier, error) {
	var a Authentifier

	if !c.IsContext(c.CurrentContext) {
		return nil, fmt.Errorf("no context found for the current context \"%s\"", c.CurrentContext)
	}
	ctx := c.Contexts[c.CurrentContext]

	switch ctx.AuthType {
	case keycloakType:
		a = &KeycloakAuthConfig{
			httpClient: httpClient,
		}
	case containerType:
		a = &ContainerAuthConfig{}
	default:
		return nil, errors.New("unknown authentication configuration")
	}

	if err := json.Unmarshal([]byte(ctx.AuthConfig), a); err != nil {
		return nil, err
	}

	return a, nil
}
