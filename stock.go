package kioom

import "net/http"

type StockBasicInfoRequest struct {
	StkCd string `json:"stk_cd"` // 거래소별 종목코드
}

type StockBasicInfoResponse struct {
	StkCd            string `json:"stk_cd"`
	StkNm            string `json:"stk_nm"`
	SetlMm           string `json:"setl_mm"`
	Fav              string `json:"fav"`        // 액면가
	Cap              string `json:"cap"`        // 자본금
	FloStk           string `json:"flo_stk"`    // 상장주식
	CrdRt            string `json:"crd_rt"`     // 신용비율
	OyrHgst          string `json:"oyr_hgst"`   // 연중최고
	OyrLwst          string `json:"oyr_lwst"`   // 연중최저
	Mac              string `json:"mac"`        // 시가총액
	MacWght          string `json:"mac_wght"`   // 시가총액비중
	ForExhRt         string `json:"for_exh_rt"` // 외인소진률
	ReplPric         string `json:"repl_pric"`  // 대용가
	Per              string `json:"per"`
	Eps              string `json:"eps"`
	Roe              string `json:"roe"`
	Pbr              string `json:"pbr"`
	Ev               string `json:"ev"`
	Bps              string `json:"bps"`
	SaleAmt          string `json:"sale_amt"`      // 매출액
	BusPro           string `json:"bus_pro"`       // 영업이익
	CupNga           string `json:"cup_nga"`       // 당기순이익
	High250          string `json:"250hgst"`       // 250최고
	Low250           string `json:"250lwst"`       // 250최저
	OpenPric         string `json:"open_pric"`     // 시가
	HighPric         string `json:"high_pric"`     // 고가
	LowPric          string `json:"low_pric"`      // 저가
	UplPric          string `json:"upl_pric"`      // 상한가
	LstPric          string `json:"lst_pric"`      // 하한가
	BasePric         string `json:"base_pric"`     // 기준가
	ExpCntrPric      string `json:"exp_cntr_pric"` // 예상체결가
	ExpCntrQty       string `json:"exp_cntr_qty"`  // 예상체결수량
	High250PricDt    string `json:"250hgst_pric_dt"`
	High250PricPreRt string `json:"250hgst_pric_pre_rt"`
	Low250PricDt     string `json:"250lwst_pric_dt"`
	Low250PricPreRt  string `json:"250lwst_pric_pre_rt"`
	CurPrc           string `json:"cur_prc"`  // 현재가
	PreSig           string `json:"pre_sig"`  // 대비기호
	PredPre          string `json:"pred_pre"` // 전일대비
	FluRt            string `json:"flu_rt"`   // 등락율
	TrdeQty          string `json:"trde_qty"` // 거래량
	TrdePre          string `json:"trde_pre"` // 거래대비
	FavUnit          string `json:"fav_unit"` // 액면가단위
	DstrStk          string `json:"dstr_stk"` // 유통주식
	DstrRt           string `json:"dstr_rt"`  // 유통비율
	ReturnCode       int    `json:"return_code"`
	ReturnMsg        string `json:"return_msg"`
}

// GetStockBasicInfo queries the basic stock info.
// API ID: ka10001
func (c *Client) GetStockBasicInfo(stockCode string) (*StockBasicInfoResponse, error) {
	reqBody := StockBasicInfoRequest{
		StkCd: stockCode,
	}

	req, err := c.newRequest(http.MethodPost, "/api/dostk/stkinfo", "ka10001", reqBody)
	if err != nil {
		return nil, err
	}

	var res StockBasicInfoResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type RealtimeItemRankRequest struct {
	QryTp string `json:"qry_tp"` // 1: 1분, 2: 10분, 3: 1시간, 4: 당일 누적, 5: 30초
}

type RealtimeItemInqRank struct {
	StkNm        string `json:"stk_nm"`
	BigdRank     string `json:"bigd_rank"`
	RankChg      string `json:"rank_chg"`
	RankChgSign  string `json:"rank_chg_sign"`
	PastCurrPrc  string `json:"past_curr_prc"`
	BaseCompSign string `json:"base_comp_sign"`
	BaseCompChgr string `json:"base_comp_chgr"`
	PrevBaseSign string `json:"prev_base_sign"`
	PrevBaseChgr string `json:"prev_base_chgr"`
	Dt           string `json:"dt"`
	Tm           string `json:"tm"`
	StkCd        string `json:"stk_cd"`
}

type RealtimeItemRankResponse struct {
	ItemInqRank []RealtimeItemInqRank `json:"item_inq_rank"`
	ReturnCode  int                   `json:"return_code"`
	ReturnMsg   string                `json:"return_msg"`
}

// GetRealtimeStockRank queries realtime top searched stocks.
// qryTp: "1" (1m), "2" (10m), "3" (1h), "4" (daily), "5" (30s)
// API ID: ka00198
func (c *Client) GetRealtimeStockRank(qryTp string) (*RealtimeItemRankResponse, error) {
	reqBody := RealtimeItemRankRequest{
		QryTp: qryTp,
	}

	req, err := c.newRequest(http.MethodPost, "/api/dostk/stkinfo", "ka00198", reqBody)
	if err != nil {
		return nil, err
	}

	var res RealtimeItemRankResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
