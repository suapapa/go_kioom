package kioom

import (
	"context"
	"net/http"
)

// USOpenOrdersRequest represents parameters for querying US stock open orders.
// Used for TR: ust21050
type USOpenOrdersRequest struct {
	OrdDt  string `json:"ord_dt"`  // Order date YYYYMMDD (empty = today)
	SlbyTp string `json:"slby_tp"` // 0: all, 1: sell, 2: buy
	StexTp string `json:"stex_tp"` // Exchange filter: NA, ND, NY (empty = all)
	StkCd  string `json:"stk_cd"`  // Ticker filter (empty = all)
}

// USOpenOrderItem represents a single open order row.
type USOpenOrderItem struct {
	OrdCntrTp   string `json:"ord_cntr_tp"` // 10: original, 11: modify, 12: cancel
	OrdNo       string `json:"ord_no"`
	OrigOrdNo   string `json:"orig_ord_no"`
	FrgnOrdID   string `json:"frgn_ord_id"`
	StexNm      string `json:"stex_nm"`
	CrncCode    string `json:"crnc_code"`
	StkCd       string `json:"stk_cd"`
	FrgnStkNm   string `json:"frgn_stk_nm"`
	FrgnTrdeTp  string `json:"frgn_trde_tp"`
	FrgnTrdeNm  string `json:"frgn_trde_nm"`
	SlbyTp      string `json:"slby_tp"`
	SlbyTpNm    string `json:"slby_tp_nm"`
	OrdQty      string `json:"ord_qty"`
	OrdUv       string `json:"ord_uv"`
	StopPric    string `json:"stop_pric"`
	CntrQty     string `json:"cntr_qty"`
	CntrUv      string `json:"cntr_uv"`
	MdfyQty     string `json:"mdfy_qty"`
	MdfyUv      string `json:"mdfy_uv"`
	CnclQty     string `json:"cncl_qty"`
	OrdRemnq    string `json:"ord_remnq"`
	OrdTime     string `json:"ord_time"`
	OrdRespTime string `json:"ord_resp_time"`
	OrdStat     string `json:"ord_stat"`
	RsrvTp      string `json:"rsrv_tp"`
	NatnNm      string `json:"natn_nm"`
}

// USOpenOrdersResponse contains US stock open order rows.
type USOpenOrdersResponse struct {
	ResultList []USOpenOrderItem `json:"result_list"`
	ReturnCode int               `json:"return_code"`
	ReturnMsg  string            `json:"return_msg"`
}

// USAccountBalanceRequest represents parameters for US stock balance inquiry.
// Used for TR: ust21070
type USAccountBalanceRequest struct {
	StexTp string `json:"stex_tp"` // Exchange: NA, ND, NY (empty = all)
	StkCd  string `json:"stk_cd"`  // Ticker filter (empty = all)
}

// USAccountBalanceItem represents a single US stock holding row.
type USAccountBalanceItem struct {
	StexNm            string `json:"stex_nm"`
	CrncCode          string `json:"crnc_code"`
	StkCd             string `json:"stk_cd"`
	FrgnStkNm         string `json:"frgn_stk_nm"`
	Qty               string `json:"qty"`
	PossQty           string `json:"poss_qty"`
	SellAlowq         string `json:"sell_alowq"`
	PredCntrSellq     string `json:"pred_cntr_sellq"`
	PredCntrBuyq      string `json:"pred_cntr_buyq"`
	TdyCntrSellq      string `json:"tdy_cntr_sellq"`
	TdyCntrBuyq       string `json:"tdy_cntr_buyq"`
	FrgnStkBookUv     string `json:"frgn_stk_book_uv"`
	NowPric           string `json:"now_pric"`
	EvltAmt           string `json:"evlt_amt"`
	PlAmt             string `json:"pl_amt"`
	PlRt              string `json:"pl_rt"`
	EvltAmtKrw        string `json:"evlt_amt_krw"`
	PlAmtKrw          string `json:"pl_amt_krw"`
	NatnNm            string `json:"natn_nm"`
	ExchRate          string `json:"exch_rate"`
	FrgnStkBookUvKrw  string `json:"frgn_stk_book_uv_krw"`
	NowPricKrw        string `json:"now_pric_krw"`
	FrgnStkBookAmt    string `json:"frgn_stk_book_amt"`
	FrgnStkBookAmtKrw string `json:"frgn_stk_book_amt_krw"`
}

// USAccountBalanceResponse contains US stock balance summary and holdings.
type USAccountBalanceResponse struct {
	StexTp        string                 `json:"stex_tp"`
	CrncCode      string                 `json:"crnc_code"`
	TotEvltAmt    string                 `json:"tot_evlt_amt"`
	TotPrchAmt    string                 `json:"tot_prch_amt"`
	TotPlAmt      string                 `json:"tot_pl_amt"`
	TotPlRt       string                 `json:"tot_pl_rt"`
	TdyBookAmt    string                 `json:"tdy_book_amt"`
	TdyPlAmt      string                 `json:"tdy_pl_amt"`
	TdyPlRt       string                 `json:"tdy_pl_rt"`
	TotEvltAmtKrw string                 `json:"tot_evlt_amt_krw"`
	TotPrchAmtKrw string                 `json:"tot_prch_amt_krw"`
	TotPlAmtKrw   string                 `json:"tot_pl_amt_krw"`
	TdyBookAmtKrw string                 `json:"tdy_book_amt_krw"`
	TdyPlAmtKrw   string                 `json:"tdy_pl_amt_krw"`
	ResultList    []USAccountBalanceItem `json:"result_list"`
	ReturnCode    int                    `json:"return_code"`
	ReturnMsg     string                 `json:"return_msg"`
}

// USDepositItem represents foreign-currency deposit details per currency.
type USDepositItem struct {
	CrncCode       string `json:"crnc_code"`
	CrncNm         string `json:"crnc_nm"`
	FcEntra        string `json:"fc_entra"`
	FcPymnAlowa    string `json:"fc_pymn_alowa"`
	FutrReplProfa  string `json:"futr_repl_profa"`
	FcBooka        string `json:"fc_booka"`
	FcOrdAlowa     string `json:"fc_ord_alowa"`
	FutrProfaBooka string `json:"futr_profa_booka"`
	FcChUncla      string `json:"fc_ch_uncla"`
	FcEtcLoana     string `json:"fc_etc_loana"`
}

// USDepositResponse contains overseas stock deposit information.
// Used for TR: ust21110
type USDepositResponse struct {
	KrwEntra   string          `json:"krw_entra"`
	ChUncla    string          `json:"ch_uncla"`
	EtcLoana   string          `json:"etc_loana"`
	ResultList []USDepositItem `json:"result_list"`
	ReturnCode int             `json:"return_code"`
	ReturnMsg  string          `json:"return_msg"`
}

// GetUSOpenOrders retrieves US stock/ETF open (unfilled) orders.
// See Kiwoom API ID: ust21050
func (c *Client) GetUSOpenOrders(ctx context.Context, req *USOpenOrdersRequest) (*USOpenOrdersResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/acnt", API_USOpenOrders, req)
	if err != nil {
		return nil, err
	}
	var res USOpenOrdersResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetUSAccountBalance retrieves US stock/ETF holdings and evaluation summary.
// See Kiwoom API ID: ust21070
func (c *Client) GetUSAccountBalance(ctx context.Context, req *USAccountBalanceRequest) (*USAccountBalanceResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/acnt", API_USAccountBalance, req)
	if err != nil {
		return nil, err
	}
	var res USAccountBalanceResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetUSDeposit retrieves overseas stock deposit balances by currency.
// See Kiwoom API ID: ust21110
func (c *Client) GetUSDeposit(ctx context.Context) (*USDepositResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/acnt", API_USDeposit, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	var res USDepositResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
