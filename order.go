package kioom

import (
	"context"
	"net/http"
)

// OrderRequest represents the payload for common stock buy/sell orders.
// Used for TR: kt10000 (Buy), kt10001 (Sell)
type OrderRequest struct {
	DmstStexTp string `json:"dmst_stex_tp"` // Domestic exchange type: KRX, NXT, SOR
	StkCd      string `json:"stk_cd"`       // Stock code
	OrdQty     string `json:"ord_qty"`      // Order quantity
	OrdUv      string `json:"ord_uv"`       // Order unit price (empty for market price)
	TrdeTp     string `json:"trde_tp"`      // Trade type (0: 보통, 3: 시장가, etc.)
	CondUv     string `json:"cond_uv"`      // Condition price (optional)
}

// OrderModifyRequest represents the payload for stock order correction.
// Used for TR: kt10002
type OrderModifyRequest struct {
	DmstStexTp string `json:"dmst_stex_tp"` // Domestic exchange type
	OrigOrdNo  string `json:"orig_ord_no"`  // Original order number
	StkCd      string `json:"stk_cd"`       // Stock code
	MdfyQty    string `json:"mdfy_qty"`     // Modified quantity
	MdfyUv     string `json:"mdfy_uv"`      // Modified unit price
	MdfyCondUv string `json:"mdfy_cond_uv"` // Modified condition price
}

// OrderCancelRequest represents the payload for stock order cancellation.
// Used for TR: kt10003
type OrderCancelRequest struct {
	DmstStexTp string `json:"dmst_stex_tp"` // Domestic exchange type
	OrigOrdNo  string `json:"orig_ord_no"`  // Original order number
	StkCd      string `json:"stk_cd"`       // Stock code
	CnclQty    string `json:"cncl_qty"`     // Cancellation quantity ('0' for all)
}

// CreditOrderRequest represents the payload for credit buy orders.
// Used for TR: kt10006
type CreditOrderRequest struct {
	OrderRequest
}

// CreditSellRequest represents the payload for credit sell orders.
// Used for TR: kt10007
type CreditSellRequest struct {
	DmstStexTp string `json:"dmst_stex_tp"`
	StkCd      string `json:"stk_cd"`
	OrdQty     string `json:"ord_qty"`
	OrdUv      string `json:"ord_uv"`
	TrdeTp     string `json:"trde_tp"`
	CrdDealTp  string `json:"crd_deal_tp"` // Credit deal type: 33 (융자), 99 (융자합)
	CrdLoanDt  string `json:"crd_loan_dt"` // Loan date (YYYYMMDD) - required for payoff
	CondUv     string `json:"cond_uv"`
}

// OrderResponse represents the common response for order requests.
type OrderResponse struct {
	OrdNo      string `json:"ord_no"`
	DmstStexTp string `json:"dmst_stex_tp"`
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

// OrderModifyResponse represents the response for order correction.
type OrderModifyResponse struct {
	OrdNo         string `json:"ord_no"`           // New order number
	BaseOrigOrdNo string `json:"base_orig_ord_no"` // Parent original order number
	MdfyQty       string `json:"mdfy_qty"`         // Modified quantity
	DmstStexTp    string `json:"dmst_stex_tp"`
	ReturnCode    int    `json:"return_code"`
	ReturnMsg     string `json:"return_msg"`
}

// OrderCancelResponse represents the response for order cancellation.
type OrderCancelResponse struct {
	OrdNo         string `json:"ord_no"`           // New order number
	BaseOrigOrdNo string `json:"base_orig_ord_no"` // Parent original order number
	CnclQty       string `json:"cncl_qty"`         // Cancelled quantity
	ReturnCode    int    `json:"return_code"`
	ReturnMsg     string `json:"return_msg"`
}

// OrderBuy sends a stock buy order.
// See Kiwoom API ID: kt10000
func (c *Client) OrderBuy(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/ordr", API_OrderBuy, req)
	if err != nil {
		return nil, err
	}
	var res OrderResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// OrderSell sends a stock sell order.
// See Kiwoom API ID: kt10001
func (c *Client) OrderSell(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/ordr", API_OrderSell, req)
	if err != nil {
		return nil, err
	}
	var res OrderResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// OrderModify sends a stock order correction request.
// See Kiwoom API ID: kt10002
func (c *Client) OrderModify(ctx context.Context, req *OrderModifyRequest) (*OrderModifyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/ordr", API_OrderModify, req)
	if err != nil {
		return nil, err
	}
	var res OrderModifyResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// OrderCancel sends a stock order cancellation request.
// See Kiwoom API ID: kt10003
func (c *Client) OrderCancel(ctx context.Context, req *OrderCancelRequest) (*OrderCancelResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/ordr", API_OrderCancel, req)
	if err != nil {
		return nil, err
	}
	var res OrderCancelResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CreditOrderBuy sends a credit buy order.
// See Kiwoom API ID: kt10006
func (c *Client) CreditOrderBuy(ctx context.Context, req *CreditOrderRequest) (*OrderResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/crdordr", API_CreditOrderBuy, req)
	if err != nil {
		return nil, err
	}
	var res OrderResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CreditOrderSell sends a credit sell order.
// See Kiwoom API ID: kt10007
func (c *Client) CreditOrderSell(ctx context.Context, req *CreditSellRequest) (*OrderResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/crdordr", API_CreditOrderSell, req)
	if err != nil {
		return nil, err
	}
	var res OrderResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CreditOrderModify sends a credit order correction request.
// See Kiwoom API ID: kt10008
func (c *Client) CreditOrderModify(ctx context.Context, req *OrderModifyRequest) (*OrderModifyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/crdordr", API_CreditModify, req)
	if err != nil {
		return nil, err
	}
	var res OrderModifyResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CreditOrderCancel sends a credit order cancellation request.
// See Kiwoom API ID: kt10009
func (c *Client) CreditOrderCancel(ctx context.Context, req *OrderCancelRequest) (*OrderCancelResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/crdordr", API_CreditCancel, req)
	if err != nil {
		return nil, err
	}
	var res OrderCancelResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
