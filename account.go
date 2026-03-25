package kioom

import "net/http"

type AccountResponse struct {
	AcctNo     string `json:"acctNo"`
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

// GetAccountNumber queries the account number belonging to the current authorized token.
// API ID: ka00001
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
