package auth

import (
	"net/http"
)

// GetToken returns current user access token.
func GetToken(httpClient *http.Client, authFileBaseName string) (string, bool, error) {
	isRefresh := false

	config, err := Load(authFileBaseName)
	if err != nil {
		return "", isRefresh, err
	}

	a, err := config.GetAuthConfig(httpClient)
	if err != nil {
		return "", isRefresh, err
	}

	if !a.Valid() {
		err := a.Renew()
		if err != nil {
			return "", isRefresh, err
		}

		err = config.UpdateContext(a)
		if err != nil {
			return "", isRefresh, err
		}

		if err := config.Save(); err != nil {
			return "", isRefresh, err
		}
		isRefresh = true
	}
	return a.GetToken(), isRefresh, nil
}
