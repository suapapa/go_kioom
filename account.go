package kioom

import "net/http"

// AccountResponse contains information about the user's trading account number.
type AccountResponse struct {
	AcctNo     string `json:"acctNo"`      // Account number
	ReturnCode int    `json:"return_code"` // 0 for success
	ReturnMsg  string `json:"return_msg"`  // Error message if return_code is not 0
}

// GetAccountNumber retrieves the account number associated with the current session.
// This requires a valid access token to be set on the client.
// See Kiwoom API ID: ka00001
func (c *Client) GetAccountNumber() (*AccountResponse, error) {
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

// DepositRequest represents the parameters for a deposit inquiry.
type DepositRequest struct {
	QryTp string `json:"qry_tp"` // 2: General, 3: Estimated
}

// DepositResponse contains the details of the account's deposit and orderable amounts.
// See Kiwoom API ID: kt00001
type DepositResponse struct {
	Deposit         string `json:"entr"`          // 예수금
	OrderAllowAmt   string `json:"ord_alow_amt"`  // 주문가능금액
	WithdrawableAmt string `json:"pymn_alow_amt"` // 출금가능금액
	D1Deposit       string `json:"d1_entra"`      // d+1추정예수금
	D2Deposit       string `json:"d2_entra"`      // d+2추정예수금
	ReturnCode      int    `json:"return_code"`   // 0 for success
	ReturnMsg       string `json:"return_msg"`    // Error message if return_code is not 0
}

// GetDeposit retrieves the deposit information for the account.
// See Kiwoom API ID: kt00001
func (c *Client) GetDeposit(req *DepositRequest) (*DepositResponse, error) {
	httpReq, err := c.newRequest(http.MethodPost, "/api/dostk/acnt", "kt00001", req)
	if err != nil {
		return nil, err
	}

	var res DepositResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// AccountBalanceRequest represents the parameters for an account balance inquiry.
type AccountBalanceRequest struct {
	QryTp     string `json:"qry_tp"`       // 1: Sum, 2: Individual
	DmstStkTP string `json:"dmst_stex_tp"` // KRX: Korea Exchange, NXT: NextTrade
}

// AccountBalanceResponse contains the summary and details of the account's stock balances.
// See Kiwoom API ID: kt00018
type AccountBalanceResponse struct {
	TotalPurchaseAmt string               `json:"tot_pur_amt"`        // 총매입금액
	TotalEvalAmt     string               `json:"tot_evlt_amt"`       // 총평가금액
	TotalEvalPL      string               `json:"tot_evlt_pl"`        // 총평가손익금액
	TotalProfitRate  string               `json:"tot_prft_rt"`        // 총수익률(%)
	AssetAmt         string               `json:"prsm_dpst_aset_amt"` // 추정예탁자산
	Items            []AccountBalanceItem `json:"acnt_evlt_remn_indv_tot"`
	ReturnCode       int                  `json:"return_code"` // 0 for success
	ReturnMsg        string               `json:"return_msg"`  // Error message if return_code is not 0
}

// AccountBalanceItem represents a single stock item in the account balance.
type AccountBalanceItem struct {
	StockCode    string `json:"stk_cd"`        // 종목번호
	StockName    string `json:"stk_nm"`        // 종목명
	EvalPL       string `json:"evltv_prft"`    // 평가손익
	ProfitRate   string `json:"prft_rt"`       // 수익률(%)
	PurchasePric string `json:"pur_pric"`      // 매입가
	Quantity     string `json:"rmnd_qty"`      // 보유수량
	TradeAbleQty string `json:"trde_able_qty"` // 매매가능수량
	CurrentPrice string `json:"cur_prc"`       // 현재가
	PurchaseAmt  string `json:"pur_amt"`       // 매입금액
	EvalAmt      string `json:"evlt_amt"`      // 평가금액
}

// GetAccountBalance retrieves the account balance information including stock items.
// See Kiwoom API ID: kt00018
func (c *Client) GetAccountBalance(req *AccountBalanceRequest) (*AccountBalanceResponse, error) {
	httpReq, err := c.newRequest(http.MethodPost, "/api/dostk/acnt", "kt00018", req)
	if err != nil {
		return nil, err
	}

	var res AccountBalanceResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
