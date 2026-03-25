package kioom

import "net/http"

// AccountResponse contains information about the user's trading account.
type AccountResponse struct {
	AcctNo     string `json:"acctNo"`      // Account number
	ReturnCode int    `json:"return_code"` // 0 for success
	ReturnMsg  string `json:"return_msg"`  // Error message if return_code is not 0
}

// GetAccountNumber retrieves the account number associated with the current session.
// This requires a valid access token to be set on the client.
// See Kiwoom API ID: ka00001
func (c *Client) GetAccountNumber() (*AccountResponse, error) {
	// body could be empty JSON object
	req, err := c.newRequest(http.MethodPost, "/api/dostk/acnt", "ka00001", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var res AccountResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
