package kioom

import (
	"fmt"
	"net/http"
)

type TokenRequest struct {
	GrantType string `json:"grant_type"`
	AppKey    string `json:"appkey"`
	SecretKey string `json:"secretkey"`
}

type TokenResponse struct {
	ExpiresDt  string `json:"expires_dt"`
	TokenType  string `json:"token_type"`
	Token      string `json:"token"`
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

type RevokeRequest struct {
	AppKey    string `json:"appkey"`
	SecretKey string `json:"secretkey"`
	Token     string `json:"token"`
}

type RevokeResponse struct {
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

// IssueToken issues a new access token and saves it in the client memory.
// API ID: au10001
func (c *Client) IssueToken() (*TokenResponse, error) {
	reqBody := TokenRequest{
		GrantType: "client_credentials",
		AppKey:    c.AppKey,
		SecretKey: c.SecretKey,
	}

	req, err := c.newRequest(http.MethodPost, "/oauth2/token", "au10001", reqBody)
	if err != nil {
		return nil, err
	}

	var res TokenResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	if res.ReturnCode != 0 {
		return &res, fmt.Errorf("issue token failed: %s", res.ReturnMsg)
	}

	// Save token so consequent calls uses this Bearer token automatically
	c.Token = res.Token
	return &res, nil
}

// RevokeToken revokes the currently issued access token.
// API ID: au10002
func (c *Client) RevokeToken() (*RevokeResponse, error) {
	reqBody := RevokeRequest{
		AppKey:    c.AppKey,
		SecretKey: c.SecretKey,
		Token:     c.Token,
	}

	req, err := c.newRequest(http.MethodPost, "/oauth2/revoke", "au10002", reqBody)
	if err != nil {
		return nil, err
	}

	var res RevokeResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	if res.ReturnCode != 0 {
		return &res, fmt.Errorf("revoke token failed: %s", res.ReturnMsg)
	}

	// Purge token logically after revocation
	c.Token = ""
	return &res, nil
}
