package main

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	kioom "github.com/suapapa/go_kioom"
)

type schemaResponse struct {
	Command  string      `json:"command"`
	Request  *typeSchema `json:"request,omitempty"`
	Response *typeSchema `json:"response,omitempty"`
	Paths    []string    `json:"paths,omitempty"`
	Notes    []string    `json:"notes,omitempty"`
}

type typeSchema struct {
	Type   string        `json:"type"`
	Fields []schemaField `json:"fields"`
}

type schemaField struct {
	Name string `json:"name"`
	JSON string `json:"json"`
	Type string `json:"type"`
}

var commandSchemaMap = map[string]struct {
	reqType  reflect.Type
	respType reflect.Type
}{
	"auth.issue":             {nil, reflect.TypeOf(kioom.TokenResponse{})},
	"auth.revoke":            {nil, reflect.TypeOf(kioom.RevokeResponse{})},
	"account.number":         {nil, reflect.TypeOf(kioom.AccountResponse{})},
	"account.deposit":        {reflect.TypeOf(kioom.DepositRequest{}), reflect.TypeOf(kioom.DepositResponse{})},
	"account.balance":        {reflect.TypeOf(kioom.AccountBalanceRequest{}), reflect.TypeOf(kioom.AccountBalanceResponse{})},
	"stock.basic":            {reflect.TypeOf(kioom.StockBasicInfoRequest{}), reflect.TypeOf(kioom.StockBasicInfoResponse{})},
	"stock.indicators":       {reflect.TypeOf(kioom.StockIndicatorRequest{}), reflect.TypeOf(kioom.StockIndicatorResponse{})},
	"stock.rank":             {reflect.TypeOf(kioom.RealtimeItemRankRequest{}), reflect.TypeOf(kioom.RealtimeItemRankResponse{})},
	"stock.volume-surge":     {reflect.TypeOf(kioom.VolumeSurgeRequest{}), reflect.TypeOf(kioom.VolumeSurgeResponse{})},
	"stock.tick-chart":       {reflect.TypeOf(kioom.StockTickChartRequest{}), reflect.TypeOf(kioom.StockTickChartResponse{})},
	"stock.minute-chart":     {reflect.TypeOf(kioom.StockMinuteChartRequest{}), reflect.TypeOf(kioom.StockMinuteChartResponse{})},
	"stock.daily-chart":      {reflect.TypeOf(kioom.StockChartRequest{}), reflect.TypeOf(kioom.StockDailyChartResponse{})},
	"stock.weekly-chart":     {reflect.TypeOf(kioom.StockChartRequest{}), reflect.TypeOf(kioom.StockWeeklyChartResponse{})},
	"stock.monthly-chart":    {reflect.TypeOf(kioom.StockChartRequest{}), reflect.TypeOf(kioom.StockMonthlyChartResponse{})},
	"stock.yearly-chart":     {reflect.TypeOf(kioom.StockChartRequest{}), reflect.TypeOf(kioom.StockYearlyChartResponse{})},
	"order.buy":              {reflect.TypeOf(kioom.OrderRequest{}), reflect.TypeOf(kioom.OrderResponse{})},
	"order.sell":             {reflect.TypeOf(kioom.OrderRequest{}), reflect.TypeOf(kioom.OrderResponse{})},
	"order.modify":           {reflect.TypeOf(kioom.OrderModifyRequest{}), reflect.TypeOf(kioom.OrderModifyResponse{})},
	"order.cancel":           {reflect.TypeOf(kioom.OrderCancelRequest{}), reflect.TypeOf(kioom.OrderCancelResponse{})},
	"us-order.buy":           {reflect.TypeOf(kioom.USOrderBuyRequest{}), reflect.TypeOf(kioom.USOrderBuyResponse{})},
	"us-order.sell":          {reflect.TypeOf(kioom.USOrderSellRequest{}), reflect.TypeOf(kioom.USOrderSellResponse{})},
	"us-order.modify":        {reflect.TypeOf(kioom.USOrderModifyRequest{}), reflect.TypeOf(kioom.USOrderModifyResponse{})},
	"us-order.cancel":        {reflect.TypeOf(kioom.USOrderCancelRequest{}), reflect.TypeOf(kioom.USOrderCancelResponse{})},
	"us-account.open-orders": {reflect.TypeOf(kioom.USOpenOrdersRequest{}), reflect.TypeOf(kioom.USOpenOrdersResponse{})},
	"us-account.balance":     {reflect.TypeOf(kioom.USAccountBalanceRequest{}), reflect.TypeOf(kioom.USAccountBalanceResponse{})},
	"us-account.deposit":     {nil, reflect.TypeOf(kioom.USDepositResponse{})},
	"api.list":               {nil, nil},
	"api.call":               {nil, nil},
}

func runSchema(args []string, stdout io.Writer, stderr io.Writer, output string) int {
	if len(args) == 0 {
		paths := make([]string, 0, len(commandSchemaMap)+len(kioom.GeneratedAPIRegistry))
		for p := range commandSchemaMap {
			paths = append(paths, p)
		}
		for _, id := range kioom.GeneratedAPIIDs() {
			paths = append(paths, "api."+id)
		}
		sort.Strings(paths)
		err := writeOK(stdout, output, schemaResponse{
			Command: "schema",
			Paths:   paths,
			Notes: []string{
				"Use `kioom schema <command-path>` to inspect request/response JSON fields.",
			},
		})
		if err != nil {
			writeErr(stderr, output, "output_error", err)
			return 1
		}
		return 0
	}

	path := strings.TrimSpace(args[0])
	if strings.HasPrefix(path, "api.") {
		return runAPISchema(path, stdout, stderr, output)
	}
	meta, ok := commandSchemaMap[path]
	if !ok {
		writeErr(stderr, output, "unknown_schema_path", fmt.Errorf("unknown command-path %q", path))
		return 2
	}

	resp := schemaResponse{Command: path}
	if meta.reqType != nil {
		reqSchema, err := toTypeSchema(meta.reqType)
		if err != nil {
			writeErr(stderr, output, "schema_error", err)
			return 1
		}
		resp.Request = reqSchema
	}
	if meta.respType != nil {
		resSchema, err := toTypeSchema(meta.respType)
		if err != nil {
			writeErr(stderr, output, "schema_error", err)
			return 1
		}
		resp.Response = resSchema
	}

	if err := writeOK(stdout, output, resp); err != nil {
		writeErr(stderr, output, "output_error", err)
		return 1
	}
	return 0
}

func toTypeSchema(t reflect.Type) (*typeSchema, error) {
	if t.Kind() != reflect.Struct {
		return nil, errors.New("schema target must be a struct")
	}
	s := &typeSchema{
		Type:   t.String(),
		Fields: make([]schemaField, 0, t.NumField()),
	}
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		tag := f.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" {
			name = f.Name
		}
		if name == "-" {
			continue
		}

		s.Fields = append(s.Fields, schemaField{
			Name: f.Name,
			JSON: name,
			Type: f.Type.String(),
		})
	}
	return s, nil
}

func runAPISchema(path string, stdout io.Writer, stderr io.Writer, output string) int {
	apiID := strings.TrimPrefix(path, "api.")
	meta, ok := kioom.GeneratedAPIRegistry[apiID]
	if !ok {
		writeErr(stderr, output, "unknown_schema_path", fmt.Errorf("unknown generated API %q", apiID))
		return 2
	}
	reqType, err := lookupGeneratedSchemaType(meta.RequestType)
	if err != nil {
		writeErr(stderr, output, "schema_error", err)
		return 1
	}
	respType, err := lookupGeneratedSchemaType(meta.ResponseType)
	if err != nil {
		writeErr(stderr, output, "schema_error", err)
		return 1
	}
	reqSchema, err := toTypeSchema(reqType)
	if err != nil {
		writeErr(stderr, output, "schema_error", err)
		return 1
	}
	resSchema, err := toTypeSchema(respType)
	if err != nil {
		writeErr(stderr, output, "schema_error", err)
		return 1
	}
	resp := schemaResponse{
		Command:  path,
		Request:  reqSchema,
		Response: resSchema,
		Notes:    []string{meta.Name, meta.Description},
	}
	if err := writeOK(stdout, output, resp); err != nil {
		writeErr(stderr, output, "output_error", err)
		return 1
	}
	return 0
}

func lookupGeneratedSchemaType(name string) (reflect.Type, error) {
	for id, meta := range kioom.GeneratedAPIRegistry {
		if meta.RequestType == name {
			req, err := kioom.NewGeneratedRequest(id)
			if err != nil {
				return nil, err
			}
			return reflect.TypeOf(req).Elem(), nil
		}
		if meta.ResponseType == name {
			res, err := kioom.NewGeneratedResponse(id)
			if err != nil {
				return nil, err
			}
			return reflect.TypeOf(res).Elem(), nil
		}
	}
	return nil, fmt.Errorf("type %q not found", name)
}
