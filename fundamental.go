package kioom

import (
	"context"
	"fmt"
	"net/http"
)

// StockForeignInvestorRequest represents the request for foreign investor trading trends.
type StockForeignInvestorRequest struct {
	StkCd string `json:"stk_cd"` // Stock code (KRX/NXT/SOR suffix supported)
}

// StockForeignInvestorItem represents one day's foreign investor activity.
type StockForeignInvestorItem struct {
	Dt             string `json:"dt"`               // Date (YYYYMMDD)
	ClosePric      string `json:"close_pric"`       // Close price
	PredPre        string `json:"pred_pre"`         // Change from previous day
	TrdeQty        string `json:"trde_qty"`         // Trading volume
	ChgQty         string `json:"chg_qty"`          // Change in holdings
	PossStkcnt     string `json:"poss_stkcnt"`      // Shares held
	Wght           string `json:"wght"`             // Ownership weight (%)
	GainPosStkcnt  string `json:"gain_pos_stkcnt"`  // Acquirable shares
	FrgnrLimit     string `json:"frgnr_limit"`      // Foreign ownership limit
	FrgnrLimitIrds string `json:"frgnr_limit_irds"` // Change in foreign limit
	LimitExhRt     string `json:"limit_exh_rt"`     // Limit exhaustion rate (%)
}

// StockForeignInvestorResponse contains daily foreign investor trading trends.
type StockForeignInvestorResponse struct {
	StkFrgnr   []StockForeignInvestorItem `json:"stk_frgnr"`
	ReturnCode int                        `json:"return_code"`
	ReturnMsg  string                     `json:"return_msg"`
}

// GetStockForeignInvestor retrieves daily foreign investor trading trends for a stock.
// See Kiwoom API ID: ka10008
func (c *Client) GetStockForeignInvestor(ctx context.Context, stockCode string) (*StockForeignInvestorResponse, error) {
	reqBody := StockForeignInvestorRequest{
		StkCd: stockCode,
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/frgnistt", API_GetStockForeignInvestor, reqBody)
	if err != nil {
		return nil, err
	}

	var res StockForeignInvestorResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// StockIndicatorRequest represents the request for comprehensive stock indicators.
type StockIndicatorRequest struct {
	StkCd string `json:"stk_cd"` // Stock code (6-digit)
}

// StockIndicatorResponse represents a unified set of queryable and calculable corporate indicators.
type StockIndicatorResponse struct {
	StkCd            string `json:"stk_cd"`               // 종목코드 (Stock code)
	StkNm            string `json:"stk_nm"`               // 종목명 (Stock name)
	Mac              string `json:"mac"`                  // 시가총액 (Market Cap)
	Per              string `json:"per"`                  // PER
	Pbr              string `json:"pbr"`                  // PBR
	Roe              string `json:"roe"`                  // ROE
	High250PricPreRt string `json:"high_250_pric_pre_rt"` // 52-week high vs current
	Low250PricPreRt  string `json:"low_250_pric_pre_rt"`  // 52-week low vs current
	ForExhRt         string `json:"for_exh_rt"`           // 외인소진률 (Foreigner exhaustion rate)
	ForQotaRt        string `json:"for_qota_rt"`          // 외국인지분율 (Foreign ownership weight from ka10008)
	SaleAmt          string `json:"sale_amt"`             // 매출액 (Revenue)
	BusPro           string `json:"bus_pro"`              // 영업이익 (Operating Profit)
	CupNga           string `json:"cup_nga"`              // 당기순이익 (Net Income)
	OperatingMargin  string `json:"operating_margin"`     // 영업이익률 (Operating profit margin, calculated)
	NetProfitMargin  string `json:"net_profit_margin"`    // 순이익률 (Net profit margin, calculated)
	ReturnCode       int    `json:"return_code"`
	ReturnMsg        string `json:"return_msg"`
}

// GetStockIndicators retrieves a comprehensive summary of stock indicators, merging basic info and foreign investor data.
func (c *Client) GetStockIndicators(ctx context.Context, stockCode string) (*StockIndicatorResponse, error) {
	basic, err := c.GetStockBasicInfo(ctx, stockCode)
	if err != nil {
		return nil, err
	}

	frgnr, err := c.GetStockForeignInvestor(ctx, stockCode)
	if err != nil {
		return nil, err
	}

	forQotaRt := ""
	if len(frgnr.StkFrgnr) > 0 {
		forQotaRt = frgnr.StkFrgnr[0].Wght
	}

	busProVal := parseFloat(basic.BusPro)
	saleAmtVal := parseFloat(basic.SaleAmt)
	opMargin := "0.00%"
	if saleAmtVal != 0 {
		opMargin = fmt.Sprintf("%.2f%%", (busProVal/saleAmtVal)*100)
	}

	cupNgaVal := parseFloat(basic.CupNga)
	npMargin := "0.00%"
	if saleAmtVal != 0 {
		npMargin = fmt.Sprintf("%.2f%%", (cupNgaVal/saleAmtVal)*100)
	}

	res := &StockIndicatorResponse{
		StkCd:            basic.StkCd,
		StkNm:            basic.StkNm,
		Mac:              basic.Mac,
		Per:              basic.Per,
		Pbr:              basic.Pbr,
		Roe:              basic.Roe,
		High250PricPreRt: basic.High250PricPreRt,
		Low250PricPreRt:  basic.Low250PricPreRt,
		ForExhRt:         basic.ForExhRt,
		ForQotaRt:        forQotaRt,
		SaleAmt:          basic.SaleAmt,
		BusPro:           basic.BusPro,
		CupNga:           basic.CupNga,
		OperatingMargin:  opMargin,
		NetProfitMargin:  npMargin,
		ReturnCode:       basic.ReturnCode,
		ReturnMsg:        basic.ReturnMsg,
	}

	return res, nil
}
