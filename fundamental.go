package kioom

import (
	"context"
	"fmt"
	"net/http"
)

// StockInstitutionsRequest represents the request to get institutional/foreign stock flows.
type StockInstitutionsRequest struct {
	StkCd string `json:"stk_cd"` // Stock code (6-digit)
}

// StockInstitutionsResponse represents the response containing institutional/foreign stock flows.
type StockInstitutionsResponse struct {
	Date             string `json:"date"`               // Date
	ClosePric        string `json:"close_pric"`         // Close price
	Pre              string `json:"pre"`                // Change from previous
	OrgnDtAcc        string `json:"orgn_dt_acc"`        // Cumulative institution volume
	OrgnDalyNettrde  string `json:"orgn_daly_nettrde"`  // Institution daily net trade
	FrgnrDalyNettrde string `json:"frgnr_daly_nettrde"` // Foreigner daily net trade
	FrgnrQotaRt      string `json:"frgnr_qota_rt"`      // Foreigner ownership ratio (외국인지분율)
	ReturnCode       int    `json:"return_code"`
	ReturnMsg        string `json:"return_msg"`
}

// GetStockInstitutions retrieves daily trading values and foreign ownership ratio for a stock.
// See Kiwoom API ID: ka10009
func (c *Client) GetStockInstitutions(ctx context.Context, stockCode string) (*StockInstitutionsResponse, error) {
	reqBody := StockInstitutionsRequest{
		StkCd: stockCode,
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/frgnistt", API_GetStockInstitutions, reqBody)
	if err != nil {
		return nil, err
	}

	var res StockInstitutionsResponse
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
	High250PricPreRt string `json:"high_250_pric_pre_rt"` // 52주일최고가대비현재가대비
	Low250PricPreRt  string `json:"low_250_pric_pre_rt"`  // 52주일최저가대비현재가대비
	ForExhRt         string `json:"for_exh_rt"`           // 외인소진률 (Foreigner exhaustion rate)
	ForQotaRt        string `json:"for_qota_rt"`          // 외국인지분율 (Foreigner ownership ratio)
	SaleAmt          string `json:"sale_amt"`             // 매출액 (Revenue)
	BusPro           string `json:"bus_pro"`              // 영업이익 (Operating Profit)
	CupNga           string `json:"cup_nga"`              // 당기순이익 (Net Income)
	OperatingMargin  string `json:"operating_margin"`     // 영업이익률 (Operating profit margin, calculated)
	NetProfitMargin  string `json:"net_profit_margin"`    // 순이익률 (Net profit margin, calculated)
	ReturnCode       int    `json:"return_code"`
	ReturnMsg        string `json:"return_msg"`
}

// GetStockIndicators retrieves a comprehensive summary of stock indicators, merging basic info and institutional data.
func (c *Client) GetStockIndicators(ctx context.Context, stockCode string) (*StockIndicatorResponse, error) {
	basic, err := c.GetStockBasicInfo(ctx, stockCode)
	if err != nil {
		return nil, err
	}

	inst, err := c.GetStockInstitutions(ctx, stockCode)
	if err != nil {
		return nil, err
	}

	// Calculate operating margin: (bus_pro / sale_amt) * 100
	busProVal := parseFloat(basic.BusPro)
	saleAmtVal := parseFloat(basic.SaleAmt)
	opMargin := "0.00%"
	if saleAmtVal != 0 {
		opMargin = fmt.Sprintf("%.2f%%", (busProVal/saleAmtVal)*100)
	}

	// Calculate net profit margin: (cup_nga / sale_amt) * 100
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
		ForQotaRt:        inst.FrgnrQotaRt,
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
