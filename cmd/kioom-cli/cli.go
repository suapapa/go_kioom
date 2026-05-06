package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	kioom "github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
	"github.com/suapapa/go_kioom/internal/kioomvalidate"
)

type envGetter func(string) string

type globalConfig struct {
	appKey    string
	secretKey string
	token     string
	mock      bool
	output    string
}

// commandFunc runs one CLI action. output is json|pretty; args are positional flags after section/action.
type commandFunc func(client *kioom.Client, output string, args []string, stdout io.Writer) error

// catalog maps "section" -> "action" -> handler. It stays aligned with [commandSchemaMap] in schema.go.
var catalog = map[string]map[string]commandFunc{
	"auth": {
		"issue":  cmdAuthIssue,
		"revoke": cmdAuthRevoke,
	},
	"account": {
		"number":  cmdAccountNumber,
		"deposit": cmdAccountDeposit,
		"balance": cmdAccountBalance,
	},
	"stock": {
		"basic":        cmdStockBasic,
		"rank":         cmdStockRank,
		"minute-chart": cmdStockMinuteChart,
	},
	"order": {
		"buy":    cmdOrderBuy,
		"sell":   cmdOrderSell,
		"modify": cmdOrderModify,
		"cancel": cmdOrderCancel,
	},
}

func run(args []string, stdout io.Writer, stderr io.Writer, getenv envGetter) int {
	cfg, rest, err := parseGlobal(args, getenv)
	if err != nil {
		writeErr(stderr, cfg.output, "invalid_arguments", err)
		return 2
	}

	if len(rest) == 0 {
		writeErr(stderr, cfg.output, "usage", errors.New(usageText()))
		return 2
	}

	if rest[0] == "schema" {
		return runSchema(rest[1:], stdout, stderr, cfg.output)
	}
	if err := requireCredentials(cfg); err != nil {
		writeErr(stderr, cfg.output, "missing_credentials", err)
		return 2
	}

	var opts []kioom.Option
	if cfg.mock {
		opts = append(opts, kioom.WithMockDomain())
	}
	client := kioom.NewClient(cfg.appKey, cfg.secretKey, opts...)
	if cfg.token != "" {
		client.SetToken(cfg.token)
	}

	if err := runCommand(client, cfg, rest, stdout); err != nil {
		writeErr(stderr, cfg.output, "command_failed", err)
		return 1
	}
	return 0
}

func parseGlobal(args []string, getenv envGetter) (globalConfig, []string, error) {
	env := kioomenv.Load(getenv)
	cfg := globalConfig{
		appKey:    env.AppKey,
		secretKey: env.SecretKey,
		token:     env.Token,
		mock:      env.Mock,
		output:    "json",
	}

	fs := flag.NewFlagSet("kioom", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.appKey, "app-key", cfg.appKey, "Kiwoom app key")
	fs.StringVar(&cfg.secretKey, "secret-key", cfg.secretKey, "Kiwoom secret key")
	fs.StringVar(&cfg.token, "token", cfg.token, "Bearer token")
	fs.BoolVar(&cfg.mock, "mock", false, "Use mock endpoint")
	fs.StringVar(&cfg.output, "output", cfg.output, "Output format: json|pretty")

	if err := fs.Parse(args); err != nil {
		return cfg, nil, err
	}

	cfg.output = strings.ToLower(strings.TrimSpace(cfg.output))
	if cfg.output != "json" && cfg.output != "pretty" {
		return cfg, nil, fmt.Errorf("unsupported --output %q", cfg.output)
	}

	return cfg, fs.Args(), nil
}

func runCommand(client *kioom.Client, cfg globalConfig, args []string, stdout io.Writer) error {
	if len(args) < 2 {
		return errors.New(usageText())
	}
	section, action, cmdArgs := args[0], args[1], args[2:]

	actions, ok := catalog[section]
	if !ok {
		return fmt.Errorf("unknown section %q", section)
	}
	h, ok := actions[action]
	if !ok {
		return fmt.Errorf("unknown %s action %q", section, action)
	}
	return h(client, cfg.output, cmdArgs, stdout)
}

func usageText() string {
	return strings.TrimSpace(`
usage:
  kioom [global-options] <section> <action> [options]
  kioom [global-options] schema <command-path>

global-options:
  --app-key, --secret-key, --token, --mock, --output

examples:
  kioom --mock auth issue
  kioom --mock --token "$KIOOM_TOKEN" account number
  kioom --mock stock basic --json '{"stk_cd":"005930"}'
  kioom schema stock.basic
`)
}

func requireCredentials(cfg globalConfig) error {
	return kioomenv.Config{AppKey: cfg.appKey, SecretKey: cfg.secretKey}.RequireAppKeys()
}

func cmdAuthIssue(c *kioom.Client, output string, _ []string, w io.Writer) error {
	res, err := c.IssueToken(context.Background())
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdAuthRevoke(c *kioom.Client, output string, _ []string, w io.Writer) error {
	res, err := c.RevokeToken(context.Background())
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdAccountNumber(c *kioom.Client, output string, _ []string, w io.Writer) error {
	res, err := c.GetAccountNumber(context.Background())
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdAccountDeposit(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.DepositRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateDepositRequest(&req); err != nil {
		return err
	}
	res, err := c.GetDeposit(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdAccountBalance(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.AccountBalanceRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateAccountBalanceRequest(&req); err != nil {
		return err
	}
	res, err := c.GetAccountBalance(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdStockBasic(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.StockBasicInfoRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateStockCode(req.StkCd); err != nil {
		return err
	}
	res, err := c.GetStockBasicInfo(context.Background(), req.StkCd)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdStockRank(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.RealtimeItemRankRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateRankRequest(&req); err != nil {
		return err
	}
	res, err := c.GetRealtimeStockRank(context.Background(), req.QryTp)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdStockMinuteChart(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.StockMinuteChartRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateMinuteChartRequest(&req); err != nil {
		return err
	}
	res, err := c.GetStockMinuteChart(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdOrderBuy(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.OrderRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateOrderRequest(&req); err != nil {
		return err
	}
	res, err := c.OrderBuy(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdOrderSell(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.OrderRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateOrderRequest(&req); err != nil {
		return err
	}
	res, err := c.OrderSell(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdOrderModify(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.OrderModifyRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateOrderModifyRequest(&req); err != nil {
		return err
	}
	res, err := c.OrderModify(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}

func cmdOrderCancel(c *kioom.Client, output string, args []string, w io.Writer) error {
	var req kioom.OrderCancelRequest
	if err := parseJSONArg(args, &req); err != nil {
		return err
	}
	if err := kioomvalidate.ValidateOrderCancelRequest(&req); err != nil {
		return err
	}
	res, err := c.OrderCancel(context.Background(), &req)
	if err != nil {
		return err
	}
	return writeOK(w, output, res)
}
