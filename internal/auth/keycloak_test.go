package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKeycloakAuthConfig(t *testing.T) {
	testCases := []struct {
		name          string
		username, url string
		wantUrl       string
	}{
		{"full_url", "username", "https://127.0.0.1:8080/auth/realms/test", "https://127.0.0.1:8080/auth/realms/test"},
		{"url_with_scheme", "username", "https://127.0.0.1:8080", "https://127.0.0.1:8080" + defaultAuthPath},
		{"url_without_scheme", "username", "127.0.0.1", defaultAuthScheme + "127.0.0.1" + defaultAuthPath},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ac := NewKeycloakAuthConfig(http.DefaultClient, tc.username, "pwd", tc.url)

			assert.Equal(t, tc.username, ac.Username)
			assert.Equal(t, tc.wantUrl, ac.BaseURL)
		})
	}
}

func TestKeyCloakLogin(t *testing.T) {
	// create a fake auth server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// handle Login() and Renew()
		if r.URL.Path == "/auth/realms/main/protocol/openid-connect/token" {
			var grantType string
			if len(r.PostForm["grant_type"]) > 0 {
				grantType = r.PostForm["grant_type"][0]
			}

			if grantType == "password" && len(r.PostForm["username"]) > 0 && r.PostForm["username"][0] == "good" {
				fmt.Fprintln(w, `{"access_token": "good_access_token", "refresh_token": "good_refresh_token"}`)
				return
			}

			if grantType == "password" && len(r.PostForm["username"]) > 0 && r.PostForm["username"][0] == "good_no_renew" {
				fmt.Fprintln(w, `{"access_token": "good_access_token", "refresh_token": "bad_refresh_token"}`)
				return
			}

			if grantType == "refresh_token" && len(r.PostForm["refresh_token"]) > 0 && r.PostForm["refresh_token"][0] == "good_refresh_token" {
				fmt.Fprintln(w, `{"access_token": "new_good_access_token", "refresh_token": "new_good_refresh_token"}`)
				return
			}

			// device authorization grant
			if grantType == "urn:ietf:params:oauth:grant-type:device_code" && len(r.PostForm["device_code"]) > 0 && r.PostForm["device_code"][0] == "good_device_code" {
				fmt.Fprintln(w, `{"access_token": "good_access_token", "refresh_token": "good_refresh_token"}`)
				return
			}
		}

		// handle device authorization grant
		if r.URL.Path == "/auth/realms/main/protocol/openid-connect/auth/device" {
			fmt.Fprintln(w, `{"device_code":"good_device_code","user_code":"good_user_code","verification_uri":"https://login.ubika.io/auth/realms/main/device","verification_uri_complete":"https://login.ubika.io/auth/realms/main/device?user_code=KWBW-IUZT","expires_in":600,"interval":1}`)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error": "internal", "error_description": "internal"}`)
	}))
	defer server.Close()

	testCases := []struct {
		name                  string
		ac                    *KeycloakAuthConfig
		token                 string
		wantErr, wantRenewErr bool
	}{
		{
			"good_login",
			NewKeycloakAuthConfig(http.DefaultClient, "good", "pwd", server.URL),
			"good_access_token", false, false,
		},
		{
			"bad_login",
			NewKeycloakAuthConfig(http.DefaultClient, "bad", "pwd", server.URL),
			"", true, true,
		},
		{
			"bad_renew",
			NewKeycloakAuthConfig(http.DefaultClient, "good_no_renew", "pwd", server.URL),
			"good_access_token", false, true,
		},
		{
			"good_device_authz_grant_flow",
			NewKeycloakAuthConfigDeviceAuthzGrant(http.DefaultClient, server.URL),
			"good_access_token", false, false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ac.Login()

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.token, tc.ac.AccessToken)

				assert.Equal(t, tc.token, tc.ac.GetToken())

				err = tc.ac.Renew()

				if tc.wantRenewErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, "new_"+tc.token, tc.ac.AccessToken)
				}
			}
		})
	}
}

func TestRequestKCToken(t *testing.T) {
	testCases := []struct {
		name       string
		path, resp string
		data       url.Values
		wantErr    error
	}{
		{"basic", "/", "simple output", url.Values{"test": {"value"}}, nil},
		{"http_error", "/wrong", "simple output", url.Values{"test": {"value"}}, &kcError{ErrorCode: "dummy"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create a fake auth server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				r.ParseForm()
				assert.True(t, assert.ObjectsAreEqualValues(r.PostForm, tc.data), "should have the same post values")
				if r.URL.Path == "/wrong" {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, `{"error": "dummy"}`)
					return
				}
				fmt.Fprintln(w, tc.resp)
			}))
			defer server.Close()

			resp, err := requestKCToken(http.DefaultClient, server.URL+tc.path, tc.data)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.resp, string(resp))
			}
		})
	}
}
