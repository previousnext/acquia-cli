package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// TokenEndpoint used to get an oauth token.
	TokenEndpoint = "https://accounts.acquia.com/api/auth/oauth/token"

	// HeaderContentType used for requests.
	HeaderContentType = "Content-Type"
	// HeaderAuthorization used for requests.
	HeaderAuthorization = "Authorization"
)

// Client for interacting with Acquia authentication.
type Client struct {
	BearerToken *BearerToken
}

// BearerToken received from the Acquia API.
type BearerToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// New client for authenticating with Acquia.
func New(key, secret string) (*Client, error) {
	bt, err := getToken(key, secret)
	if err != nil {
		return nil, err
	}

	return &Client{bt}, nil
}

// WrapRequest adds the Authorization header for Acquia API requests.
func (c *Client) WrapRequest(r *http.Request) {
	r.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.BearerToken.AccessToken))
}

// ExchangeKeysForBearerToken exchanges an Acquia API key and secret for an
// oauth2 bearer token.
func getToken(key, secret string) (*BearerToken, error) {
	auth := fmt.Sprintf("%s:%s", key, secret)
	authEnc := base64.StdEncoding.EncodeToString([]byte(auth))
	body := []byte("grant_type=client_credentials&scope=")

	req, err := http.NewRequest("POST", TokenEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(HeaderContentType, "application/x-www-form-urlencoded")
	req.Header.Set(HeaderAuthorization, fmt.Sprintf("Basic %s", authEnc))

	client := &http.Client{}

	resp, err :=  client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := &BearerToken{}
	err = json.Unmarshal(data, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
