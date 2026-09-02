package kioom

import (
	"context"
	"net/http"
)

// AccountResponse contains information about the user's trading account number.
type AccountResponse struct {
	AcctNo     string `json:"acctNo"`      // Account number
	ReturnCode int    `json:"return_code"` // 0 for success
	ReturnMsg  string `json:"return_msg"`  // Error message if return_code is not 0
}

// GetAccountNumber retrieves the account number associated with the current session.
// This requires a valid access token to be set on the client.
// See Kiwoom API ID: ka00001
func (c *Client) GetAccountNumber(ctx context.Context) (*AccountResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/acnt", API_GetAccountNumber, map[string]interface{}{})
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
func (c *Client) GetDeposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/acnt", API_GetDeposit, req)
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
	TotalLoanAmt     string               `json:"tot_loan_amt"`       // 총대출금
	TotalCrdLoanAmt  string               `json:"tot_crd_loan_amt"`   // 총융자금액
	TotalCrdLsAmt    string               `json:"tot_crd_ls_amt"`     // 총대주금액
	Items            []AccountBalanceItem `json:"acnt_evlt_remn_indv_tot"`
	ReturnCode       int                  `json:"return_code"` // 0 for success
	ReturnMsg        string               `json:"return_msg"`  // Error message if return_code is not 0
}

// AccountBalanceItem represents a single stock item in the account balance.
type AccountBalanceItem struct {
	StockCode     string `json:"stk_cd"`          // 종목번호
	StockName     string `json:"stk_nm"`          // 종목명
	EvalPL        string `json:"evltv_prft"`      // 평가손익
	ProfitRate    string `json:"prft_rt"`         // 수익률(%)
	PurchasePric  string `json:"pur_pric"`        // 매입가
	PredClosePric string `json:"pred_close_pric"` // 전일종가
	Quantity      string `json:"rmnd_qty"`        // 보유수량
	TradeAbleQty  string `json:"trde_able_qty"`   // 매매가능수량
	CurrentPrice  string `json:"cur_prc"`         // 현재가
	PredBuyQty    string `json:"pred_buyq"`       // 전일매수수량
	PredSellQty   string `json:"pred_sellq"`      // 전일매도수량
	TdyBuyQty     string `json:"tdy_buyq"`        // 금일매수수량
	TdySellQty    string `json:"tdy_sellq"`       // 금일매도수량
	PurchaseAmt   string `json:"pur_amt"`         // 매입금액
	PurchaseCmsn  string `json:"pur_cmsn"`        // 매입수수료
	EvalAmt       string `json:"evlt_amt"`        // 평가금액
	SellCmsn      string `json:"sell_cmsn"`       // 평가수수료
	Tax           string `json:"tax"`             // 세금
	SumCmsn       string `json:"sum_cmsn"`        // 수수료합
	PossRt        string `json:"poss_rt"`         // 보유비중(%)
	CrdTp         string `json:"crd_tp"`          // 신용구분
	CrdTpNm       string `json:"crd_tp_nm"`       // 신용구분명
	CrdLoanDt     string `json:"crd_loan_dt"`     // 대출일
}

// GetAccountBalance retrieves the account balance information including stock items.
// See Kiwoom API ID: kt00018
func (c *Client) GetAccountBalance(ctx context.Context, req *AccountBalanceRequest) (*AccountBalanceResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/acnt", API_GetAccountBalance, req)
	if err != nil {
		return nil, err
	}

	var res AccountBalanceResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
