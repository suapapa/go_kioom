package kioom

import (
	"context"
	"net/http"
)

// StockMinuteChartRequest payload for querying minute chart data.
type StockMinuteChartRequest struct {
	StkCd      string `json:"stk_cd"`       // Stock code
	TicScope   string `json:"tic_scope"`    // 1:1m, 3:3m, 5:5m, 10:10m, 15:15m, 30:30m, 45:45m, 60:60m
	UpdStkpcTp string `json:"upd_stkpc_tp"` // 0 or 1
	BaseDt     string `json:"base_dt"`      // YYYYMMDD (optional)
}

// StockMinutePole represents a single minute candle's data.
type StockMinutePole struct {
	CurPrc     string `json:"cur_prc"`      // Close price
	TrdeQty    string `json:"trde_qty"`     // Trading volume
	CntrTm     string `json:"cntr_tm"`      // Contract time (YYYYMMDDHHMMSS)
	OpenPric   string `json:"open_pric"`    // Open price
	HighPric   string `json:"high_pric"`    // High price
	LowPric    string `json:"low_pric"`     // Low price
	AccTrdeQty string `json:"acc_trde_qty"` // Accumulated trading volume
	PredPre    string `json:"pred_pre"`     // Change from previous day
	PredPreSig string `json:"pred_pre_sig"` // Change sign (1: Up limit, 2: Up, 3: Even, 4: Down limit, 5: Down)
}

// StockMinuteChartResponse contains a list of minute candles.
type StockMinuteChartResponse struct {
	StkCd              string            `json:"stk_cd"`
	StkMinPoleChartQry []StockMinutePole `json:"stk_min_pole_chart_qry"`
	ReturnCode         int               `json:"return_code"`
	ReturnMsg          string            `json:"return_msg"`
}

// GetStockMinuteChart retrieves minute-level chart data for a specific stock.
// ticScope: "1" (1m), "3" (3m), "5" (5m), "10" (10m), etc.
// updStkpcTp: "0" (Standard) or "1" (Adjusted).
// baseDt: Target date in YYYYMMDD format.
// See Kiwoom API ID: ka10080
func (c *Client) GetStockMinuteChart(ctx context.Context, req *StockMinuteChartRequest) (*StockMinuteChartResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/chart", API_GetMinuteChart, req)
	if err != nil {
		return nil, err
	}

	var res StockMinuteChartResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// StockTickChartRequest payload for querying tick chart data.
type StockTickChartRequest struct {
	StkCd      string `json:"stk_cd"`       // Stock code (6-digit)
	TicScope   string `json:"tic_scope"`    // 1:1tic, 3:3tic, 5:5tic, 10:10tic, 30:30tic
	UpdStkpcTp string `json:"upd_stkpc_tp"` // 0: Standard, 1: Adjusted
}

// StockTickPole represents a single tick candle's data.
type StockTickPole struct {
	CurPrc     string `json:"cur_prc"`      // Close price (현재가)
	TrdeQty    string `json:"trde_qty"`     // Trading volume (거래량)
	CntrTm     string `json:"cntr_tm"`      // Contract time (체결시간, YYYYMMDDHHMMSS)
	OpenPric   string `json:"open_pric"`    // Open price (시가)
	HighPric   string `json:"high_pric"`    // High price (고가)
	LowPric    string `json:"low_pric"`     // Low price (저가)
	PredPre    string `json:"pred_pre"`     // Change from previous day (전일대비)
	PredPreSig string `json:"pred_pre_sig"` // Change sign (1: Up limit, 2: Up, 3: Even, 4: Down limit, 5: Down)
}

// StockTickChartResponse contains a list of tick candles.
type StockTickChartResponse struct {
	StkCd          string          `json:"stk_cd"`
	LastTicCnt     string          `json:"last_tic_cnt"`
	StkTicChartQry []StockTickPole `json:"stk_tic_chart_qry"`
	ReturnCode     int             `json:"return_code"`
	ReturnMsg      string          `json:"return_msg"`
}

// GetStockTickChart retrieves tick-level chart data for a specific stock.
// See Kiwoom API ID: ka10079
func (c *Client) GetStockTickChart(ctx context.Context, req *StockTickChartRequest) (*StockTickChartResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/chart", API_GetTickChart, req)
	if err != nil {
		return nil, err
	}

	var res StockTickChartResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// StockChartRequest payload for querying daily, weekly, monthly, and yearly chart data.
type StockChartRequest struct {
	StkCd      string `json:"stk_cd"`       // Stock code (6-digit)
	BaseDt     string `json:"base_dt"`      // Base date (YYYYMMDD)
	UpdStkpcTp string `json:"upd_stkpc_tp"` // 0: Standard, 1: Adjusted
}

// StockChartPole represents a single candle's data for daily, weekly, monthly, and yearly charts.
type StockChartPole struct {
	CurPrc     string `json:"cur_prc"`                // Close price (현재가)
	TrdeQty    string `json:"trde_qty"`               // Trading volume (거래량)
	TrdePrica  string `json:"trde_prica"`             // Trading value (거래대금)
	Dt         string `json:"dt"`                     // Date (일자 YYYYMMDD)
	OpenPric   string `json:"open_pric"`              // Open price (시가)
	HighPric   string `json:"high_pric"`              // High price (고가)
	LowPric    string `json:"low_pric"`               // Low price (저가)
	PredPre    string `json:"pred_pre,omitempty"`     // Change from previous day (전일대비, not present in yearly)
	PredPreSig string `json:"pred_pre_sig,omitempty"` // Change sign (전일대비기호, not present in yearly)
	TrdeTernRt string `json:"trde_tern_rt,omitempty"` // Transaction turn rate (거래회전율, not present in yearly)
}

// StockDailyChartResponse contains daily chart data.
type StockDailyChartResponse struct {
	StkCd             string           `json:"stk_cd"`
	StkDtPoleChartQry []StockChartPole `json:"stk_dt_pole_chart_qry"`
	ReturnCode        int              `json:"return_code"`
	ReturnMsg         string           `json:"return_msg"`
}

// StockWeeklyChartResponse contains weekly chart data.
type StockWeeklyChartResponse struct {
	StkCd              string           `json:"stk_cd"`
	StkStkPoleChartQry []StockChartPole `json:"stk_stk_pole_chart_qry"`
	ReturnCode         int              `json:"return_code"`
	ReturnMsg          string           `json:"return_msg"`
}

// StockMonthlyChartResponse contains monthly chart data.
type StockMonthlyChartResponse struct {
	StkCd              string           `json:"stk_cd"`
	StkMthPoleChartQry []StockChartPole `json:"stk_mth_pole_chart_qry"`
	ReturnCode         int              `json:"return_code"`
	ReturnMsg          string           `json:"return_msg"`
}

// StockYearlyChartResponse contains yearly chart data.
type StockYearlyChartResponse struct {
	StkCd             string           `json:"stk_cd"`
	StkYrPoleChartQry []StockChartPole `json:"stk_yr_pole_chart_qry"`
	ReturnCode        int              `json:"return_code"`
	ReturnMsg         string           `json:"return_msg"`
}

// GetStockDailyChart retrieves daily chart data for a specific stock.
// See Kiwoom API ID: ka10081
func (c *Client) GetStockDailyChart(ctx context.Context, req *StockChartRequest) (*StockDailyChartResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/chart", API_GetDailyChart, req)
	if err != nil {
		return nil, err
	}

	var res StockDailyChartResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetStockWeeklyChart retrieves weekly chart data for a specific stock.
// See Kiwoom API ID: ka10082
func (c *Client) GetStockWeeklyChart(ctx context.Context, req *StockChartRequest) (*StockWeeklyChartResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/chart", API_GetWeeklyChart, req)
	if err != nil {
		return nil, err
	}

	var res StockWeeklyChartResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetStockMonthlyChart retrieves monthly chart data for a specific stock.
// See Kiwoom API ID: ka10083
func (c *Client) GetStockMonthlyChart(ctx context.Context, req *StockChartRequest) (*StockMonthlyChartResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/chart", API_GetMonthlyChart, req)
	if err != nil {
		return nil, err
	}

	var res StockMonthlyChartResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetStockYearlyChart retrieves yearly chart data for a specific stock.
// See Kiwoom API ID: ka10094
func (c *Client) GetStockYearlyChart(ctx context.Context, req *StockChartRequest) (*StockYearlyChartResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/chart", API_GetYearlyChart, req)
	if err != nil {
		return nil, err
	}

	var res StockYearlyChartResponse
	if err := c.do(httpReq, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
