// Package mcpkioom implements the Kiwoom MCP tool surface on top of [kioom.Client].
// Transports ([mcp.StdioTransport], [mcp.NewSSEHandler], etc.) attach to the same
// constructed [mcp.Server] so stdio and HTTP stay in sync.
package mcpkioom

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	kioom "github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomvalidate"
)

const defaultServerName = "kioom-mcp"

// NewServer registers Kiwoom tools on a new [mcp.Server]. Callers run it with any
// [mcp.Transport] (e.g. [mcp.StdioTransport]) or expose it via [mcp.NewSSEHandler].
func NewServer(client *kioom.Client, impl *mcp.Implementation, opts *mcp.ServerOptions) *mcp.Server {
	if impl == nil {
		impl = &mcp.Implementation{Name: defaultServerName, Version: "0.1.0"}
	}
	servOpts := opts
	if servOpts == nil {
		servOpts = &mcp.ServerOptions{Instructions: serverInstructions()}
	} else if servOpts.Instructions == "" {
		clone := *servOpts
		clone.Instructions = serverInstructions()
		servOpts = &clone
	}

	s := mcp.NewServer(impl, servOpts)
	addTools(s, client)
	return s
}

func serverInstructions() string {
	return `Kiwoom Open API MCP server. Set KIOOM_APP_KEY and KIOOM_SECRET_KEY (and optionally KIOOM_TOKEN, KIOOM_MOCK=true) like kioom-cli. Tools mirror kioom-cli sections auth, account, stock, and order.`
}

func addTools(s *mcp.Server, c *kioom.Client) {
	type noArgs struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_issue",
		Description: "Issue a new OAuth access token (Kiwoom au10001). Token is stored on the client for subsequent calls.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ noArgs) (*mcp.CallToolResult, *kioom.TokenResponse, error) {
		res, err := c.IssueToken()
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_revoke",
		Description: "Revoke the current access token (Kiwoom au10002).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ noArgs) (*mcp.CallToolResult, *kioom.RevokeResponse, error) {
		res, err := c.RevokeToken()
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "account_number",
		Description: "Get the account number.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ noArgs) (*mcp.CallToolResult, *kioom.AccountResponse, error) {
		res, err := c.GetAccountNumber()
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "account_deposit",
		Description: "Query deposit information; body matches kioom.DepositRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.DepositRequest) (*mcp.CallToolResult, *kioom.DepositResponse, error) {
		if err := kioomvalidate.ValidateDepositRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetDeposit(&in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "account_balance",
		Description: "Query account balance; body matches kioom.AccountBalanceRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.AccountBalanceRequest) (*mcp.CallToolResult, *kioom.AccountBalanceResponse, error) {
		if err := kioomvalidate.ValidateAccountBalanceRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetAccountBalance(&in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_basic",
		Description: "Stock basic info for stk_cd (six digits).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in struct {
		StkCd string `json:"stk_cd" jsonschema:"required six-digit stock code"`
	}) (*mcp.CallToolResult, *kioom.StockBasicInfoResponse, error) {
		if err := kioomvalidate.ValidateStockCode(in.StkCd); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockBasicInfo(in.StkCd)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_rank",
		Description: "Realtime stock rank; body matches kioom.RealtimeItemRankRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.RealtimeItemRankRequest) (*mcp.CallToolResult, *kioom.RealtimeItemRankResponse, error) {
		if err := kioomvalidate.ValidateRankRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetRealtimeStockRank(in.QryTp)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_minute_chart",
		Description: "Minute chart; body matches kioom.StockMinuteChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockMinuteChartRequest) (*mcp.CallToolResult, *kioom.StockMinuteChartResponse, error) {
		if err := kioomvalidate.ValidateMinuteChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockMinuteChart(&in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_buy",
		Description: "Place a buy order; body matches kioom.OrderRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderRequest) (*mcp.CallToolResult, *kioom.OrderResponse, error) {
		if err := kioomvalidate.ValidateOrderRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderBuy(&in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_sell",
		Description: "Place a sell order; body matches kioom.OrderRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderRequest) (*mcp.CallToolResult, *kioom.OrderResponse, error) {
		if err := kioomvalidate.ValidateOrderRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderSell(&in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_modify",
		Description: "Modify an order; body matches kioom.OrderModifyRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderModifyRequest) (*mcp.CallToolResult, *kioom.OrderModifyResponse, error) {
		if err := kioomvalidate.ValidateOrderModifyRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderModify(&in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_cancel",
		Description: "Cancel an order; body matches kioom.OrderCancelRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderCancelRequest) (*mcp.CallToolResult, *kioom.OrderCancelResponse, error) {
		if err := kioomvalidate.ValidateOrderCancelRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderCancel(&in)
		return nil, res, err
	})
}
