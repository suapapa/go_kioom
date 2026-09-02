package kioom

import (
	"context"
	"net/http"
)

// USOrderBuyRequest represents the payload for US stock/ETF buy orders.
// Used for TR: ust20000
type USOrderBuyRequest struct {
	StexTp string `json:"stex_tp"` // Exchange: NA (AMEX), ND (NASDAQ), NY (NYSE)
	StkCd  string `json:"stk_cd"`  // Ticker symbol (e.g. NVDA, SPY)
	OrdQty string `json:"ord_qty"` // Order quantity
	OrdUv  string `json:"ord_uv"`  // Order unit price (empty for market orders)
	TrdeTp string `json:"trde_tp"` // Trade type: 00 limit, 03 market, 26 VWAP limit, etc.
}

// USOrderSellRequest represents the payload for US stock/ETF sell orders.
// Used for TR: ust20001
type USOrderSellRequest struct {
	StkCd    string `json:"stk_cd"`    // Ticker symbol
	StexTp   string `json:"stex_tp"`   // Exchange: NA, ND, NY
	OrdQty   string `json:"ord_qty"`   // Order quantity
	OrdUv    string `json:"ord_uv"`    // Order unit price (empty for market orders)
	StopPric string `json:"stop_pric"` // STOP price (required for trde_tp 34/35)
	TrdeTp   string `json:"trde_tp"`   // Trade type: 00 limit, 03 market, 34 STOP LIMIT, 35 STOP, etc.
}

// USOrderModifyRequest represents the payload for US stock/ETF order correction.
// Used for TR: ust20002
type USOrderModifyRequest struct {
	OrigOrdNo string `json:"orig_ord_no"` // Original order number from order response
	StexTp    string `json:"stex_tp"`     // Exchange: NA, ND, NY
	StkCd     string `json:"stk_cd"`      // Ticker symbol
	MdfyUv    string `json:"mdfy_uv"`     // Modified unit price
	StopPric  string `json:"stop_pric"`   // STOP price (required when original order is STOP type)
}

// USOrderCancelRequest represents the payload for US stock/ETF order cancellation.
// Used for TR: ust20003
type USOrderCancelRequest struct {
	OrigOrdNo string `json:"orig_ord_no"` // Original order number
	StexTp    string `json:"stex_tp"`     // Exchange: NA, ND, NY
	StkCd     string `json:"stk_cd"`      // Ticker symbol
}

// USOrderBuyResponse represents the response for US buy orders.
type USOrderBuyResponse struct {
	AcntNm         string `json:"acnt_nm"`
	StkNm          string `json:"stk_nm"`
	OrdNo          string `json:"ord_no"`
	FcEntra        string `json:"fc_entra"`        // Foreign currency deposit
	TdyRebuyUseda  string `json:"tdy_rebuy_useda"` // Today's rebuy used amount
	PredRebuyUseda string `json:"pred_rebuy_useda"`
	TrstProfCh     string `json:"trst_prof_ch"` // Margin used
	ReturnCode     int    `json:"return_code"`
	ReturnMsg      string `json:"return_msg"`
}

// USOrderSellResponse represents the response for US sell orders.
type USOrderSellResponse struct {
	AcntNm         string `json:"acnt_nm"`
	StkNm          string `json:"stk_nm"`
	OrdNo          string `json:"ord_no"`
	PossQty        string `json:"poss_qty"`        // Holdings quantity
	TdyReselUsedq  string `json:"tdy_resel_usedq"` // Today's resell used quantity
	PredReselUsedq string `json:"pred_resel_usedq"`
	ReturnCode     int    `json:"return_code"`
	ReturnMsg      string `json:"return_msg"`
}

// USOrderModifyResponse represents the response for US order correction.
type USOrderModifyResponse struct {
	AcntNm         string `json:"acnt_nm"`
	StkNm          string `json:"stk_nm"`
	OrdNo          string `json:"ord_no"`
	FcEntra        string `json:"fc_entra"`
	TdyRebuyUseda  string `json:"tdy_rebuy_useda"`
	PredRebuyUseda string `json:"pred_rebuy_useda"`
	TrstProfCh     string `json:"trst_prof_ch"`
	MdfyOrdQty     string `json:"mdfy_ord_qty"`
	ReturnCode     int    `json:"return_code"`
	ReturnMsg      string `json:"return_msg"`
}

// USOrderCancelResponse represents the response for US order cancellation.
type USOrderCancelResponse struct {
	AcntNm     string `json:"acnt_nm"`
	StkNm      string `json:"stk_nm"`
	OrdNo      string `json:"ord_no"`
	CnclOrdQty string `json:"cncl_ord_qty"`
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

// USOrderBuy sends a US stock/ETF buy order.
// See Kiwoom API ID: ust20000
func (c *Client) USOrderBuy(ctx context.Context, req *USOrderBuyRequest) (*USOrderBuyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/ordr", API_USOrderBuy, req)
	if err != nil {
		return nil, err
	}
	var res USOrderBuyResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// USOrderSell sends a US stock/ETF sell order.
// See Kiwoom API ID: ust20001
func (c *Client) USOrderSell(ctx context.Context, req *USOrderSellRequest) (*USOrderSellResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/ordr", API_USOrderSell, req)
	if err != nil {
		return nil, err
	}
	var res USOrderSellResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// USOrderModify sends a US stock/ETF order correction request.
// See Kiwoom API ID: ust20002
func (c *Client) USOrderModify(ctx context.Context, req *USOrderModifyRequest) (*USOrderModifyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/ordr", API_USOrderModify, req)
	if err != nil {
		return nil, err
	}
	var res USOrderModifyResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// USOrderCancel sends a US stock/ETF order cancellation request.
// See Kiwoom API ID: ust20003
func (c *Client) USOrderCancel(ctx context.Context, req *USOrderCancelRequest) (*USOrderCancelResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/us/ordr", API_USOrderCancel, req)
	if err != nil {
		return nil, err
	}
	var res USOrderCancelResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
