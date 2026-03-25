package kioom

import (
	"fmt"
	"net/http"
)

// TokenRequest contains the credentials required to request a new access token.
type TokenRequest struct {
	GrantType string `json:"grant_type"`
	AppKey    string `json:"appkey"`
	SecretKey string `json:"secretkey"`
}

// TokenResponse contains the issued access token and its metadata.
type TokenResponse struct {
	ExpiresDt  string `json:"expires_dt"`  // Format: YYYY-MM-DD HH:MM:SS
	TokenType  string `json:"token_type"`  // Should be "Bearer"
	Token      string `json:"token"`       // The access token itself
	ReturnCode int    `json:"return_code"` // 0 for success
	ReturnMsg  string `json:"return_msg"`  // Error message if return_code is not 0
}

// RevokeRequest contains the credentials required to revoke an access token.
type RevokeRequest struct {
	AppKey    string `json:"appkey"`
	SecretKey string `json:"secretkey"`
	Token     string `json:"token"`
}

// RevokeResponse indicates whether the revocation was successful.
type RevokeResponse struct {
	ReturnCode int    `json:"return_code"` // 0 for success
	ReturnMsg  string `json:"return_msg"`  // Error message if return_code is not 0
}

// IssueToken issues a new access token using the client's credentials.
// It automatically saves the issued token to the Client instance.
// See Kiwoom API ID: au10001
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
// After successfully revoking the token, the Client's internal Token field is cleared.
// See Kiwoom API ID: au10002
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
