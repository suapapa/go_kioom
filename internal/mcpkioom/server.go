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
	return `Kiwoom Open API MCP server. Requires KIOOM_APP_KEY and KIOOM_SECRET_KEY environment variables (and optionally KIOOM_TOKEN, KIOOM_MOCK=true for testing). Provides tools for authentication (auth_*), account management (account_*), stock querying (stock_*), and trading/ordering (order_*).`
}

func addTools(s *mcp.Server, c *kioom.Client) {
	type noArgs struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_issue",
		Description: "Issue a new OAuth access token (Kiwoom au10001). Token is stored on the client for subsequent calls.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ noArgs) (*mcp.CallToolResult, *kioom.TokenResponse, error) {
		res, err := c.IssueToken(ctx)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_revoke",
		Description: "Revoke the current access token (Kiwoom au10002).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ noArgs) (*mcp.CallToolResult, *kioom.RevokeResponse, error) {
		res, err := c.RevokeToken(ctx)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "account_number",
		Description: "Get the account number.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ noArgs) (*mcp.CallToolResult, *kioom.AccountResponse, error) {
		res, err := c.GetAccountNumber(ctx)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "account_deposit",
		Description: "Query deposit information; body matches kioom.DepositRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.DepositRequest) (*mcp.CallToolResult, *kioom.DepositResponse, error) {
		if err := kioomvalidate.ValidateDepositRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetDeposit(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "account_balance",
		Description: "Query account balance; body matches kioom.AccountBalanceRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.AccountBalanceRequest) (*mcp.CallToolResult, *kioom.AccountBalanceResponse, error) {
		if err := kioomvalidate.ValidateAccountBalanceRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetAccountBalance(ctx, &in)
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
		res, err := c.GetStockBasicInfo(ctx, in.StkCd)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_indicators",
		Description: "Comprehensive fundamental stock indicators for stk_cd (six digits) including computed margins and foreign ownership weight (ka10008 wght).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in struct {
		StkCd string `json:"stk_cd" jsonschema:"required six-digit stock code"`
	}) (*mcp.CallToolResult, *kioom.StockIndicatorResponse, error) {
		if err := kioomvalidate.ValidateStockCode(in.StkCd); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockIndicators(ctx, in.StkCd)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_rank",
		Description: "Realtime stock rank; body matches kioom.RealtimeItemRankRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.RealtimeItemRankRequest) (*mcp.CallToolResult, *kioom.RealtimeItemRankResponse, error) {
		if err := kioomvalidate.ValidateRankRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetRealtimeStockRank(ctx, in.QryTp)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_volume_surge",
		Description: "Query stock trading volume surge or drop; body matches kioom.VolumeSurgeRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.VolumeSurgeRequest) (*mcp.CallToolResult, *kioom.VolumeSurgeResponse, error) {
		if err := kioomvalidate.ValidateVolumeSurgeRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetVolumeSurge(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_tick_chart",
		Description: "Tick chart; body matches kioom.StockTickChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockTickChartRequest) (*mcp.CallToolResult, *kioom.StockTickChartResponse, error) {
		if err := kioomvalidate.ValidateTickChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockTickChart(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_minute_chart",
		Description: "Minute chart; body matches kioom.StockMinuteChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockMinuteChartRequest) (*mcp.CallToolResult, *kioom.StockMinuteChartResponse, error) {
		if err := kioomvalidate.ValidateMinuteChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockMinuteChart(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_daily_chart",
		Description: "Daily chart; body matches kioom.StockChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockChartRequest) (*mcp.CallToolResult, *kioom.StockDailyChartResponse, error) {
		if err := kioomvalidate.ValidateChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockDailyChart(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_weekly_chart",
		Description: "Weekly chart; body matches kioom.StockChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockChartRequest) (*mcp.CallToolResult, *kioom.StockWeeklyChartResponse, error) {
		if err := kioomvalidate.ValidateChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockWeeklyChart(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_monthly_chart",
		Description: "Monthly chart; body matches kioom.StockChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockChartRequest) (*mcp.CallToolResult, *kioom.StockMonthlyChartResponse, error) {
		if err := kioomvalidate.ValidateChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockMonthlyChart(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "stock_yearly_chart",
		Description: "Yearly chart; body matches kioom.StockChartRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.StockChartRequest) (*mcp.CallToolResult, *kioom.StockYearlyChartResponse, error) {
		if err := kioomvalidate.ValidateChartRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.GetStockYearlyChart(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_buy",
		Description: "Place a buy order; body matches kioom.OrderRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderRequest) (*mcp.CallToolResult, *kioom.OrderResponse, error) {
		if err := kioomvalidate.ValidateOrderRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderBuy(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_sell",
		Description: "Place a sell order; body matches kioom.OrderRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderRequest) (*mcp.CallToolResult, *kioom.OrderResponse, error) {
		if err := kioomvalidate.ValidateOrderRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderSell(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_modify",
		Description: "Modify an order; body matches kioom.OrderModifyRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderModifyRequest) (*mcp.CallToolResult, *kioom.OrderModifyResponse, error) {
		if err := kioomvalidate.ValidateOrderModifyRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderModify(ctx, &in)
		return nil, res, err
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "order_cancel",
		Description: "Cancel an order; body matches kioom.OrderCancelRequest.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in kioom.OrderCancelRequest) (*mcp.CallToolResult, *kioom.OrderCancelResponse, error) {
		if err := kioomvalidate.ValidateOrderCancelRequest(&in); err != nil {
			return nil, nil, err
		}
		res, err := c.OrderCancel(ctx, &in)
		return nil, res, err
	})
}
