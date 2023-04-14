package auth

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// creates and load config for testing
	c := newConfig()
	c.path = filepath.Join(t.TempDir(), ".appsecctl")
	err := c.load()

	assert.NoError(t, err)
	assert.Equal(t, map[string]Context{}, c.Contexts)
}

func TestUseContext(t *testing.T) {
	testCases := []struct {
		name    string
		ctxName string
	}{
		{"use existing context", "default"},
		{"use unexisting context will create it", "ctx_name"},
	}

	c := genTestConfig(t)
	c.UseContext("default")

	for _, tc := range testCases {
		c.UseContext(tc.ctxName)

		assert.Equal(t, tc.ctxName, c.CurrentContext)

		_, inMap := c.Contexts[tc.ctxName]
		assert.True(t, inMap)
	}
}

func TestSave(t *testing.T) {
	c := genTestConfig(t)

	// adding new empty context to config
	c.UseContext("titi")
	err := c.Save()

	assert.NoError(t, err)

	// read file to check if new context exists
	err = c.load()
	assert.NoError(t, err)
	assert.Equal(t, "titi", c.CurrentContext)

	_, inMap := c.Contexts["titi"]
	assert.True(t, inMap)
}

func TestUpdateContext(t *testing.T) {
	testCases := []struct {
		name    string
		ctxName string
	}{
		{"current context not exists", "default"},
		{"current context exists", "ctx_name"},
	}

	for _, tc := range testCases {
		c := initConfig(t)
		if tc.name == "current context exists" {
			c.CurrentContext = tc.ctxName
		}
		newAuth := KeycloakAuthConfig{}

		err := c.UpdateContext(&newAuth)

		assert.NoError(t, err)
		assert.Equal(t, tc.ctxName, c.CurrentContext)

		_, inMap := c.Contexts[tc.ctxName]
		assert.True(t, inMap)
	}
}

func TestIsContext(t *testing.T) {
	testCases := []struct {
		name     string
		ctxName  string
		expected bool
	}{
		{"context exists", "default", true},
		{"context not exists", "toto", false},
	}

	c := initConfig(t)
	c.Contexts["default"] = Context{}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, c.IsContext(tc.ctxName))
	}
}

func TestGetAuthConfig(t *testing.T) {
	testCases := []struct {
		name     string
		ctxName  string
		authType Type
	}{
		{"use toto context", "toto", containerType},
		{"switch to default context", "default", keycloakType},
	}

	c := genTestConfig(t)

	for _, tc := range testCases {
		c.UseContext(tc.ctxName)
		a, err := c.GetAuthConfig(http.DefaultClient)
		assert.NoError(t, err)
		assert.Equal(t, tc.authType, a.GetType())
	}
}

// genTestConfig generates a fake config for testing.
func genTestConfig(t *testing.T) Config {
	c := initConfig(t)
	a1 := KeycloakAuthConfig{}
	a2 := ContainerAuthConfig{}

	c.UseContext("default")
	c.UpdateContext(&a1)

	c.UseContext("toto")
	c.UpdateContext(&a2)

	_ = c.Save()

	return c
}

// initConfig init empty config for testing.
func initConfig(t *testing.T) Config {
	c := newConfig()
	c.path = filepath.Join(t.TempDir(), "appsecctl")

	err := c.Save()
	assert.NoError(t, err)

	err = c.load()
	assert.NoError(t, err)

	return c
}
