package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"

	kioom "github.com/suapapa/go_kioom"
)

var (
	reStockCode = regexp.MustCompile(`^\d{6}$`)
)

func parseJSONArg(args []string, out interface{}) error {
	fs := flag.NewFlagSet("json", flag.ContinueOnError)
	fs.SetOutput(bytes.NewBuffer(nil))

	jsonStr := ""
	fs.StringVar(&jsonStr, "json", "", "Request body in JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(jsonStr) == "" {
		return errors.New("--json is required")
	}

	dec := json.NewDecoder(strings.NewReader(jsonStr))
	dec.DisallowUnknownFields()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("invalid --json payload: %w", err)
	}
	if dec.More() {
		return errors.New("invalid --json payload: multiple JSON values")
	}
	return nil
}

func validateStockCode(code string) error {
	if !reStockCode.MatchString(strings.TrimSpace(code)) {
		return fmt.Errorf("invalid stock code %q: expected 6 digits", code)
	}
	return nil
}

func validateDepositRequest(req *kioom.DepositRequest) error {
	switch req.QryTp {
	case "2", "3":
		return nil
	default:
		return fmt.Errorf("invalid qry_tp %q: expected one of [2,3]", req.QryTp)
	}
}

func validateAccountBalanceRequest(req *kioom.AccountBalanceRequest) error {
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

func validateRankRequest(req *kioom.RealtimeItemRankRequest) error {
	switch req.QryTp {
	case "1", "2", "3", "4", "5":
		return nil
	default:
		return fmt.Errorf("invalid qry_tp %q: expected one of [1,2,3,4,5]", req.QryTp)
	}
}

func validateMinuteChartRequest(req *kioom.StockMinuteChartRequest) error {
	if err := validateStockCode(req.StkCd); err != nil {
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

func validateOrderRequest(req *kioom.OrderRequest) error {
	if err := validateStockCode(req.StkCd); err != nil {
		return err
	}
	switch req.DmstStexTp {
	case "KRX", "NXT", "SOR":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT,SOR]", req.DmstStexTp)
	}
}

func validateOrderModifyRequest(req *kioom.OrderModifyRequest) error {
	if err := validateStockCode(req.StkCd); err != nil {
		return err
	}
	switch req.DmstStexTp {
	case "KRX", "NXT", "SOR":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT,SOR]", req.DmstStexTp)
	}
}

func validateOrderCancelRequest(req *kioom.OrderCancelRequest) error {
	if err := validateStockCode(req.StkCd); err != nil {
		return err
	}
	switch req.DmstStexTp {
	case "KRX", "NXT", "SOR":
		return nil
	default:
		return fmt.Errorf("invalid dmst_stex_tp %q: expected one of [KRX,NXT,SOR]", req.DmstStexTp)
	}
}
