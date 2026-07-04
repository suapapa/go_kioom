package kioom

import (
	"strings"
	"time"
)

const tokenExpiryBuffer = 60 * time.Second

var kstLocation = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.FixedZone("KST", 9*60*60)
	}
	return loc
}()

// parseExpiresDt parses Kiwoom token expiry timestamps.
// The API returns compact timestamps such as "20241107083713".
func parseExpiresDt(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		"20060102150405",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, kstLocation); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errInvalidExpiresDt
}

var errInvalidExpiresDt = &expiresDtError{}

type expiresDtError struct{}

func (e *expiresDtError) Error() string { return "invalid expires_dt format" }

func (c *Client) tokenValidAt(now time.Time) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tokenValidLocked(now)
}

func (c *Client) tokenValidLocked(now time.Time) bool {
	if c.token == "" {
		return false
	}
	if c.tokenExpiresAt.IsZero() {
		// Tokens from KIOOM_TOKEN or SetToken have unknown expiry; keep using
		// them until the API reports an auth error.
		return true
	}
	return now.Add(tokenExpiryBuffer).Before(c.tokenExpiresAt)
}

func (c *Client) cachedTokenResponse() *TokenResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expiresDt := ""
	if !c.tokenExpiresAt.IsZero() {
		expiresDt = c.tokenExpiresAt.In(kstLocation).Format("20060102150405")
	}

	return &TokenResponse{
		ExpiresDt:  expiresDt,
		TokenType:  "Bearer",
		Token:      c.token,
		ReturnCode: 0,
		ReturnMsg:  "success (cached)",
	}
}

func (c *Client) clearToken() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = ""
	c.tokenExpiresAt = time.Time{}
}

func (c *Client) setTokenWithExpiry(token string, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
	c.tokenExpiresAt = expiresAt
}

type apiResponseMeta struct {
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

func isTokenAuthError(meta apiResponseMeta) bool {
	if meta.ReturnCode == 0 {
		return false
	}
	// Kiwoom reports invalid/expired tokens as return_code 3 (e.g. "8005:Token이 유효하지 않습니다").
	if meta.ReturnCode == 3 {
		return true
	}
	msg := meta.ReturnMsg
	return strings.Contains(msg, "8005") ||
		(strings.Contains(msg, "Token") && strings.Contains(msg, "유효"))
}

func isAuthEndpoint(path string) bool {
	return path == "/oauth2/token" || path == "/oauth2/revoke"
}
