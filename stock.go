package kioom

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// StockBasicInfoRequest represents the payload for querying stock metadata.
type StockBasicInfoRequest struct {
	StkCd string `json:"stk_cd"` // Stock code from the exchange
}

// StockBasicInfoResponse contains detailed information about a single stock.
type StockBasicInfoResponse struct {
	StkCd            string `json:"stk_cd"`        // Stock code
	StkNm            string `json:"stk_nm"`        // Stock name
	SetlMm           string `json:"setl_mm"`       // Settlement month
	Fav              string `json:"fav"`           // Par value
	Cap              string `json:"cap"`           // Capital
	FloStk           string `json:"flo_stk"`       // Number of floating shares
	CrdRt            string `json:"crd_rt"`        // Credit ratio
	OyrHgst          string `json:"oyr_hgst"`      // Year-to-date high
	OyrLwst          string `json:"oyr_lwst"`      // Year-to-date low
	Mac              string `json:"mac"`           // Market cap
	MacWght          string `json:"mac_wght"`      // Market cap weight
	ForExhRt         string `json:"for_exh_rt"`    // Foreigner exhaustion rate
	ReplPric         string `json:"repl_pric"`     // Collateral value
	Per              string `json:"per"`           // Price-to-earnings ratio
	Eps              string `json:"eps"`           // Earnings per share
	Roe              string `json:"roe"`           // Return on equity
	Pbr              string `json:"pbr"`           // Price-to-book ratio
	Ev               string `json:"ev"`            // Enterprise value
	Bps              string `json:"bps"`           // Book value per share
	SaleAmt          string `json:"sale_amt"`      // Revenue/Sales amount
	BusPro           string `json:"bus_pro"`       // Operating profit
	CupNga           string `json:"cup_nga"`       // Net income
	High250          string `json:"250hgst"`       // 250-day high
	Low250           string `json:"250lwst"`       // 250-day low
	OpenPric         string `json:"open_pric"`     // Opening price
	HighPric         string `json:"high_pric"`     // High price
	LowPric          string `json:"low_pric"`      // Low price
	UplPric          string `json:"upl_pric"`      // Upper limit price
	LstPric          string `json:"lst_pric"`      // Lower limit price
	BasePric         string `json:"base_pric"`     // Base price
	ExpCntrPric      string `json:"exp_cntr_pric"` // Expected contract price
	ExpCntrQty       string `json:"exp_cntr_qty"`  // Expected contract volume
	High250PricDt    string `json:"250hgst_pric_dt"`
	High250PricPreRt string `json:"250hgst_pric_pre_rt"`
	Low250PricDt     string `json:"250lwst_pric_dt"`
	Low250PricPreRt  string `json:"250lwst_pric_pre_rt"`
	CurPrc           string `json:"cur_prc"`  // Current price
	PreSig           string `json:"pre_sig"`  // Change sign
	PredPre          string `json:"pred_pre"` // Change from previous day
	FluRt            string `json:"flu_rt"`   // Fluctuation rate
	TrdeQty          string `json:"trde_qty"` // Trading volume
	TrdePre          string `json:"trde_pre"` // Comparison to previous volume
	FavUnit          string `json:"fav_unit"` // Par value unit
	DstrStk          string `json:"dstr_stk"` // Number of distributed shares
	DstrRt           string `json:"dstr_rt"`  // Distribution ratio
	ReturnCode       int    `json:"return_code"`
	ReturnMsg        string `json:"return_msg"`
}

// GetStockBasicInfo retrieves comprehensive metadata for a specific stock.
// See Kiwoom API ID: ka10001
func (c *Client) GetStockBasicInfo(ctx context.Context, stockCode string) (*StockBasicInfoResponse, error) {
	reqBody := StockBasicInfoRequest{
		StkCd: stockCode,
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/stkinfo", API_GetStockBasicInfo, reqBody)
	if err != nil {
		return nil, err
	}

	var res StockBasicInfoResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// RealtimeItemRankRequest payload for querying searched stock rankings.
type RealtimeItemRankRequest struct {
	QryTp string `json:"qry_tp"` // 1: 1m, 2: 10m, 3: 1h, 4: Daily, 5: 30s
}

// RealtimeItemInqRank represents a single stock's ranking information.
type RealtimeItemInqRank struct {
	StkNm        string `json:"stk_nm"`         // Stock name
	BigdRank     string `json:"bigd_rank"`      // Rank
	RankChg      string `json:"rank_chg"`       // Change in rank
	RankChgSign  string `json:"rank_chg_sign"`  // Sign of rank change
	PastCurrPrc  string `json:"past_curr_prc"`  // Previous price
	BaseCompSign string `json:"base_comp_sign"` // Comparison sign
	BaseCompChgr string `json:"base_comp_chgr"` // Change amount
	PrevBaseSign string `json:"prev_base_sign"` // Previous comparison sign
	PrevBaseChgr string `json:"prev_base_chgr"` // Previous change amount
	Dt           string `json:"dt"`             // Date
	Tm           string `json:"tm"`             // Time
	StkCd        string `json:"stk_cd"`         // Stock code
}

// RealtimeItemRankResponse contains a list of ranked stocks.
type RealtimeItemRankResponse struct {
	ItemInqRank []RealtimeItemInqRank `json:"item_inq_rank"`
	ReturnCode  int                   `json:"return_code"`
	ReturnMsg   string                `json:"return_msg"`
}

// GetRealtimeStockRank retrieves the most searched stocks in realtime.
// Use qryTp to specify the time window: "1" (1m), "2" (10m), "3" (1h), "4" (daily), "5" (30s)
// See Kiwoom API ID: ka00198
func (c *Client) GetRealtimeStockRank(ctx context.Context, qryTp string) (*RealtimeItemRankResponse, error) {
	reqBody := RealtimeItemRankRequest{
		QryTp: qryTp,
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/stkinfo", API_GetRealtimeRank, reqBody)
	if err != nil {
		return nil, err
	}

	var res RealtimeItemRankResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// parseFloat converts a string to float64 safely, stripping spaces, commas, and leading signs.
func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	if s[0] == '+' {
		s = s[1:]
	}
	s = strings.ReplaceAll(s, ",", "")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// VolumeSurgeRequest represents the request payload for ka10023.
type VolumeSurgeRequest struct {
	MrktTp    string `json:"mrkt_tp"`     // 000:전체, 001:코스피, 101:코스닥
	SortTp    string `json:"sort_tp"`     // 1:급증량, 2:급증률, 3:급감량, 4:급감률
	TmTp      string `json:"tm_tp"`       // 1:분, 2:전일
	TrdeQtyTp string `json:"trde_qty_tp"` // 5:5천주이상, 10:만주이상, 50:5만주이상, 100:10만주이상, 200:20만주이상, 300:30만주이상, 500:50만주이상, 1000:백만주이상
	Tm        string `json:"tm"`          // 분 입력
	StkCnd    string `json:"stk_cnd"`     // 0:전체조회, 1:관리종목제외, 3:우선주제외, 11:정리매매종목제외...
	PricTp    string `json:"pric_tp"`     // 0:전체조회, 2:5만원이상, 5:1만원이상, 6:5천원이상, 8:1천원이상, 9:10만원이상
	StexTp    string `json:"stex_tp"`     // 1:KRX, 2:NXT, 3:통합
}

// VolumeSurgeItem represents a single item in the volume surge list.
type VolumeSurgeItem struct {
	StkCd       string `json:"stk_cd"`        // 종목코드
	StkNm       string `json:"stk_nm"`        // 종목명
	CurPrc      string `json:"cur_prc"`       // 현재가
	PredPreSig  string `json:"pred_pre_sig"`  // 전일대비기호
	PredPre     string `json:"pred_pre"`      // 전일대비
	FluRt       string `json:"flu_rt"`        // 등락률
	PrevTrdeQty string `json:"prev_trde_qty"` // 이전거래량
	NowTrdeQty  string `json:"now_trde_qty"`  // 현재거래량
	SdninQty    string `json:"sdnin_qty"`     // 급증량
	SdninRt     string `json:"sdnin_rt"`      // 급증률
}

// VolumeSurgeResponse represents the response containing volume surge list.
type VolumeSurgeResponse struct {
	TrdeQtySdnin []VolumeSurgeItem `json:"trde_qty_sdnin"`
	ReturnCode   int               `json:"return_code"`
	ReturnMsg    string            `json:"return_msg"`
}

// GetVolumeSurge queries stocks with sudden volume surge or drop (ka10023).
func (c *Client) GetVolumeSurge(ctx context.Context, reqBody *VolumeSurgeRequest) (*VolumeSurgeResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/dostk/rkinfo", API_GetVolumeSurge, reqBody)
	if err != nil {
		return nil, err
	}

	var res VolumeSurgeResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
