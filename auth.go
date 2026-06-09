package kioom

import (
	"context"
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
func (c *Client) IssueToken(ctx context.Context) (*TokenResponse, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// Double-check if the token was already refreshed by another concurrent request.
	if currentToken := c.Token(); currentToken != "" {
		return &TokenResponse{
			Token:      currentToken,
			ReturnCode: 0,
			ReturnMsg:  "success (cached)",
		}, nil
	}

	reqBody := TokenRequest{
		GrantType: "client_credentials",
		AppKey:    c.appKey,
		SecretKey: c.secretKey,
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/oauth2/token", API_IssueToken, reqBody)
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
	c.SetToken(res.Token)
	return &res, nil
}

// IssueTokenForce forcibly issues a new token even if a cached one exists.
func (c *Client) IssueTokenForce(ctx context.Context) (*TokenResponse, error) {
	c.mu.Lock()
	c.token = ""
	c.mu.Unlock()
	return c.IssueToken(ctx)
}

// RevokeToken revokes the currently issued access token.
// After successfully revoking the token, the Client's internal Token field is cleared.
// See Kiwoom API ID: au10002
func (c *Client) RevokeToken(ctx context.Context) (*RevokeResponse, error) {
	reqBody := RevokeRequest{
		AppKey:    c.appKey,
		SecretKey: c.secretKey,
		Token:     c.Token(),
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/oauth2/revoke", API_RevokeToken, reqBody)
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
	c.SetToken("")
	return &res, nil
}
