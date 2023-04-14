package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// authResponse is a Keycloak authentication response.
type authResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type KeycloakDeviceAuthzResp struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURLComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// KeycloakAuthConfig define authentication configuration.
type KeycloakAuthConfig struct {
	httpClient       *http.Client
	Type             Type   `json:"type"`
	Username         string `json:"username"`
	BaseURL          string `json:"base_url"`
	password         string `json:"-"`
	deviceAuthzGrant bool   `json:"-"`

	authResponse
}

var (
	// defaultAuthPath is the authentication server default path.
	defaultAuthPath string = "/auth/realms/main"

	// defaultAuthScheme is the default scheme for auth server.
	defaultAuthScheme string = "https://"
)

func NewKeycloakAuthConfig(client *http.Client, username, pwd, url string) *KeycloakAuthConfig {
	ac := KeycloakAuthConfig{
		httpClient: client,
		Type:       keycloakType,
		Username:   username,
		password:   pwd,
	}
	ac.setURL(url)
	return &ac
}

func NewKeycloakAuthConfigDeviceAuthzGrant(client *http.Client, url string) *KeycloakAuthConfig {
	ac := KeycloakAuthConfig{
		httpClient:       client,
		Type:             keycloakType,
		deviceAuthzGrant: true,
	}
	ac.setURL(url)
	return &ac
}

// Login authenticates user with the authentication server.
func (ac *KeycloakAuthConfig) Login() error {
	if ac.deviceAuthzGrant {
		return ac.loginDeviceAuthzGrant()
	}
	return ac.loginWithPassword()
}

func (ac *KeycloakAuthConfig) loginDeviceAuthzGrant() error {
	resp, err := requestKCToken(
		ac.httpClient,
		ac.BaseURL+"/protocol/openid-connect/auth/device",
		url.Values{
			"client_id": {"appsecctl"},
			"scope":     {"offline_access"},
		},
	)
	if err != nil {
		return err
	}

	da := KeycloakDeviceAuthzResp{}
	err = json.Unmarshal(resp, &da)
	if err != nil {
		return err
	}
	if da.Interval == 0 {
		da.Interval = 5
	}
	if da.ExpiresIn == 0 {
		da.ExpiresIn = 300
	}

	fmt.Printf("To log in, open the page %s with a web browser and enter the code %s to authenticate.\n", da.VerificationURI, da.UserCode)

	t := time.NewTicker(time.Duration(da.Interval) * time.Second)
	defer t.Stop()

LOOP:
	for {
		select {
		case <-time.After(time.Duration(da.ExpiresIn) * time.Second):
			return fmt.Errorf("Too slow! Authentication code has expired. :(")
		case <-t.C:
			resp, err = requestKCToken(
				ac.httpClient,
				ac.BaseURL+"/protocol/openid-connect/token",
				url.Values{
					"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
					"client_id":   {"appsecctl"},
					"device_code": {da.DeviceCode},
				},
			)
			if err != nil {
				if e, ok := err.(*kcError); ok && (e.ErrorCode == "authorization_pending" || e.ErrorCode == "slow_down") {
					continue
				}
				return err
			}

			break LOOP
		}
	}

	return json.Unmarshal(resp, &ac)
}

func (ac *KeycloakAuthConfig) loginWithPassword() error {
	resp, err := requestKCToken(
		ac.httpClient,
		ac.BaseURL+"/protocol/openid-connect/token",
		url.Values{
			"username":   {ac.Username},
			"password":   {ac.password},
			"client_id":  {"appsecctl"},
			"grant_type": {"password"},
			"scope":      {"offline_access"},
		},
	)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, &ac)
}

// Renew renews AccessToken.
func (ac *KeycloakAuthConfig) Renew() error {
	resp, err := requestKCToken(
		ac.httpClient,
		ac.BaseURL+"/protocol/openid-connect/token",
		url.Values{
			"client_id":     {"appsecctl"},
			"grant_type":    {"refresh_token"},
			"refresh_token": {ac.RefreshToken},
		},
	)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, &ac)
}

// GetToken returns an access token for the current user.
func (ac *KeycloakAuthConfig) GetToken() string {
	return ac.AccessToken
}

// Valid returns true if current AccessToken is valid or not expired.
func (ac *KeycloakAuthConfig) Valid() bool {
	_, keyFunc := NewRSAAuth()
	token, err := jwt.Parse(ac.AccessToken, keyFunc)

	// let the backend validate the token
	// force renew if expired
	if token.Valid {
		return true
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) == 0
	} else {
		return false
	}
}

// GetType returns the type of the Auth Configuration.
func (ac *KeycloakAuthConfig) GetType() Type {
	return keycloakType
}

func (ac *KeycloakAuthConfig) setURL(url string) {
	// add scheme if missing in the URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = defaultAuthScheme + url
	}

	// add base path to the URL if missing
	if !strings.Contains(url, "auth/realms/") {
		url += defaultAuthPath
	}

	ac.BaseURL = strings.TrimSuffix(url, "/")
}

// requestKCToken requests new token from Keycloak authentication server.
func requestKCToken(client *http.Client, url string, data url.Values) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	rawBody = bytes.TrimSuffix(rawBody, []byte("\n"))

	if resp.StatusCode != 200 {
		if err != nil {
			return nil, err
		}
		return nil, parseError(rawBody)
	}

	return rawBody, err
}

type kcError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *kcError) Error() string {
	return e.ErrorDescription
}

func parseError(data []byte) error {
	e := &kcError{}
	err := json.Unmarshal(data, e)
	if err != nil {
		return err
	}
	return e
}
