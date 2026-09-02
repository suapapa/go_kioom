// Package kioomvalidate contains request validation shared by the CLI, MCP server, and other tools.
package kioomvalidate

import (
	"fmt"
	"regexp"
	"strings"

	kioom "github.com/suapapa/go_kioom"
)

var reStockCode = regexp.MustCompile(`^\d{6}$`)

// ValidateStockCode checks that code is a six-digit stock symbol.
func ValidateStockCode(code string) error {
	if !reStockCode.MatchString(strings.TrimSpace(code)) {
		return fmt.Errorf("invalid stock code %q: expected 6 digits", code)
	}
	return nil
}

// ValidateDepositRequest validates qry_tp for deposit queries.
func ValidateDepositRequest(req *kioom.DepositRequest) error {
	switch req.QryTp {
	case "2", "3":
		return nil
	default:
		return fmt.Errorf("invalid qry_tp %q: expected one of [2,3]", req.QryTp)
	}
}

// ValidateAccountBalanceRequest validates domestic balance query parameters.
func ValidateAccountBalanceRequest(req *kioom.AccountBalanceRequest) error {
	switch req.QryTp {
	case "1", "2":
	default:
		return fmt.Errorf("invalid qry_tp %q: expected one of [1,2]", req.QryTp)
	}
	switch req.DmstStkTP {
	case "KRX", "NXT":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT]", req.DmstStkTP)
	}
}

// ValidateRankRequest validates realtime rank query type.
func ValidateRankRequest(req *kioom.RealtimeItemRankRequest) error {
	switch req.QryTp {
	case "1", "2", "3", "4", "5":
		return nil
	default:
		return fmt.Errorf("invalid qry_tp %q: expected one of [1,2,3,4,5]", req.QryTp)
	}
}

// ValidateMinuteChartRequest validates minute chart parameters.
func ValidateMinuteChartRequest(req *kioom.StockMinuteChartRequest) error {
	if err := ValidateStockCode(req.StkCd); err != nil {
		return err
	}

	switch req.TicScope {
	case "1", "3", "5", "10", "15", "30", "45", "60":
	default:
		return fmt.Errorf("invalid tic_scope %q", req.TicScope)
	}

	switch req.UpdStkpcTp {
	case "0", "1":
	default:
		return fmt.Errorf("invalid upd_stkpc_tp %q: expected one of [0,1]", req.UpdStkpcTp)
	}

	return nil
}

// ValidateTickChartRequest validates tick chart parameters.
func ValidateTickChartRequest(req *kioom.StockTickChartRequest) error {
	if err := ValidateStockCode(req.StkCd); err != nil {
		return err
	}

	switch req.TicScope {
	case "1", "3", "5", "10", "30":
	default:
		return fmt.Errorf("invalid tic_scope %q: expected one of [1,3,5,10,30]", req.TicScope)
	}

	switch req.UpdStkpcTp {
	case "0", "1":
	default:
		return fmt.Errorf("invalid upd_stkpc_tp %q: expected one of [0,1]", req.UpdStkpcTp)
	}

	return nil
}

// ValidateChartRequest validates daily, weekly, monthly, and yearly chart parameters.
func ValidateChartRequest(req *kioom.StockChartRequest) error {
	if err := ValidateStockCode(req.StkCd); err != nil {
		return err
	}

	if len(req.BaseDt) != 8 {
		return fmt.Errorf("invalid base_dt %q: expected 8 characters (YYYYMMDD)", req.BaseDt)
	}

	switch req.UpdStkpcTp {
	case "0", "1":
	default:
		return fmt.Errorf("invalid upd_stkpc_tp %q: expected one of [0,1]", req.UpdStkpcTp)
	}

	return nil
}

// ValidateOrderRequest validates a new order request.
func ValidateOrderRequest(req *kioom.OrderRequest) error {
	if err := ValidateStockCode(req.StkCd); err != nil {
		return err
	}
	switch req.DmstStexTp {
	case "KRX", "NXT", "SOR":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT,SOR]", req.DmstStexTp)
	}
}

// ValidateOrderModifyRequest validates an order modification request.
func ValidateOrderModifyRequest(req *kioom.OrderModifyRequest) error {
	if err := ValidateStockCode(req.StkCd); err != nil {
		return err
	}
	switch req.DmstStexTp {
	case "KRX", "NXT", "SOR":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT,SOR]", req.DmstStexTp)
	}
}

// ValidateOrderCancelRequest validates an order cancel request.
func ValidateOrderCancelRequest(req *kioom.OrderCancelRequest) error {
	if err := ValidateStockCode(req.StkCd); err != nil {
		return err
	}
	switch req.DmstStexTp {
	case "KRX", "NXT", "SOR":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT,SOR]", req.DmstStexTp)
	}
}

// ValidateVolumeSurgeRequest validates volume surge request parameters.
func ValidateVolumeSurgeRequest(req *kioom.VolumeSurgeRequest) error {
	switch req.MrktTp {
	case "000", "001", "101":
	default:
		return fmt.Errorf("invalid mrkt_tp %q: expected one of [000, 001, 101]", req.MrktTp)
	}

	switch req.SortTp {
	case "1", "2", "3", "4":
	default:
		return fmt.Errorf("invalid sort_tp %q: expected one of [1, 2, 3, 4]", req.SortTp)
	}

	switch req.TmTp {
	case "1", "2":
	default:
		return fmt.Errorf("invalid tm_tp %q: expected one of [1, 2]", req.TmTp)
	}

	switch req.TrdeQtyTp {
	case "5", "10", "50", "100", "200", "300", "500", "1000":
	default:
		return fmt.Errorf("invalid trde_qty_tp %q: expected one of [5, 10, 50, 100, 200, 300, 500, 1000]", req.TrdeQtyTp)
	}

	switch req.StexTp {
	case "1", "2", "3":
	default:
		return fmt.Errorf("invalid stex_tp %q: expected one of [1, 2, 3]", req.StexTp)
	}

	return nil
}

var reUSTicker = regexp.MustCompile(`^[A-Za-z0-9]{1,12}$`)

// ValidateUSTicker checks that code is a valid US stock/ETF ticker (1-12 alphanumeric).
func ValidateUSTicker(code string) error {
	code = strings.TrimSpace(code)
	if !reUSTicker.MatchString(code) {
		return fmt.Errorf("invalid US ticker %q: expected 1-12 alphanumeric characters", code)
	}
	return nil
}

// ValidateUSExchangeType checks the US exchange type code.
func ValidateUSExchangeType(stexTp string) error {
	if err := validateUSExchangeTypeOptional(stexTp, false); err != nil {
		return err
	}
	return nil
}

func validateUSExchangeTypeOptional(stexTp string, allowEmpty bool) error {
	stexTp = strings.TrimSpace(stexTp)
	if stexTp == "" {
		if allowEmpty {
			return nil
		}
		return fmt.Errorf("invalid stex_tp %q: expected one of [NA,ND,NY]", stexTp)
	}
	switch stexTp {
	case "NA", "ND", "NY":
		return nil
	default:
		return fmt.Errorf("invalid stex_tp %q: expected one of [NA,ND,NY]", stexTp)
	}
}

var usBuyTradeTypes = map[string]struct{}{
	"00": {}, "03": {}, "26": {}, "27": {}, "30": {}, "36": {}, "37": {},
}

var usSellTradeTypes = map[string]struct{}{
	"00": {}, "03": {}, "26": {}, "27": {}, "30": {}, "33": {}, "34": {}, "35": {}, "36": {}, "37": {},
}

// ValidateUSOrderBuyRequest validates a US buy order request.
func ValidateUSOrderBuyRequest(req *kioom.USOrderBuyRequest) error {
	if err := ValidateUSExchangeType(req.StexTp); err != nil {
		return err
	}
	if err := ValidateUSTicker(req.StkCd); err != nil {
		return err
	}
	if _, ok := usBuyTradeTypes[req.TrdeTp]; !ok {
		return fmt.Errorf("invalid trde_tp %q: expected one of [00,03,26,27,30,36,37]", req.TrdeTp)
	}
	return nil
}

// ValidateUSOrderSellRequest validates a US sell order request.
func ValidateUSOrderSellRequest(req *kioom.USOrderSellRequest) error {
	if err := ValidateUSExchangeType(req.StexTp); err != nil {
		return err
	}
	if err := ValidateUSTicker(req.StkCd); err != nil {
		return err
	}
	if _, ok := usSellTradeTypes[req.TrdeTp]; !ok {
		return fmt.Errorf("invalid trde_tp %q: expected one of [00,03,26,27,30,33,34,35,36,37]", req.TrdeTp)
	}
	return nil
}

// ValidateUSOrderModifyRequest validates a US order modification request.
func ValidateUSOrderModifyRequest(req *kioom.USOrderModifyRequest) error {
	if err := ValidateUSExchangeType(req.StexTp); err != nil {
		return err
	}
	if err := ValidateUSTicker(req.StkCd); err != nil {
		return err
	}
	if strings.TrimSpace(req.OrigOrdNo) == "" {
		return fmt.Errorf("orig_ord_no is required")
	}
	return nil
}

// ValidateUSOrderCancelRequest validates a US order cancel request.
func ValidateUSOrderCancelRequest(req *kioom.USOrderCancelRequest) error {
	if err := ValidateUSExchangeType(req.StexTp); err != nil {
		return err
	}
	if err := ValidateUSTicker(req.StkCd); err != nil {
		return err
	}
	if strings.TrimSpace(req.OrigOrdNo) == "" {
		return fmt.Errorf("orig_ord_no is required")
	}
	return nil
}

// ValidateUSOpenOrdersRequest validates US open orders query parameters.
func ValidateUSOpenOrdersRequest(req *kioom.USOpenOrdersRequest) error {
	if err := validateUSExchangeTypeOptional(req.StexTp, true); err != nil {
		return err
	}
	if req.StkCd != "" {
		if err := ValidateUSTicker(req.StkCd); err != nil {
			return err
		}
	}
	if req.OrdDt != "" && len(req.OrdDt) != 8 {
		return fmt.Errorf("invalid ord_dt %q: expected 8 characters (YYYYMMDD)", req.OrdDt)
	}
	switch req.SlbyTp {
	case "", "0", "1", "2":
		return nil
	default:
		return fmt.Errorf("invalid slby_tp %q: expected one of [0,1,2]", req.SlbyTp)
	}
}

// ValidateUSAccountBalanceRequest validates US balance query parameters.
func ValidateUSAccountBalanceRequest(req *kioom.USAccountBalanceRequest) error {
	if err := validateUSExchangeTypeOptional(req.StexTp, true); err != nil {
		return err
	}
	if req.StkCd != "" {
		return ValidateUSTicker(req.StkCd)
	}
	return nil
}
