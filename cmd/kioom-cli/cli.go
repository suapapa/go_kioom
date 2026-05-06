package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	kioom "github.com/suapapa/go_kioom"
)

type envGetter func(string) string

type globalConfig struct {
	appKey    string
	secretKey string
	token     string
	mock      bool
	output    string
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

	client := kioom.NewClient(cfg.appKey, cfg.secretKey, cfg.mock)
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
	cfg := globalConfig{
		appKey:    getenv("KIOOM_APP_KEY"),
		secretKey: getenv("KIOOM_SECRET_KEY"),
		token:     getenv("KIOOM_TOKEN"),
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

	section := args[0]
	action := args[1]
	cmdArgs := args[2:]

	switch section {
	case "auth":
		return runAuth(client, cfg.output, action, cmdArgs, stdout)
	case "account":
		return runAccount(client, cfg.output, action, cmdArgs, stdout)
	case "stock":
		return runStock(client, cfg.output, action, cmdArgs, stdout)
	case "order":
		return runOrder(client, cfg.output, action, cmdArgs, stdout)
	default:
		return fmt.Errorf("unknown section %q", section)
	}
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
	if strings.TrimSpace(cfg.appKey) == "" || strings.TrimSpace(cfg.secretKey) == "" {
		return errors.New("app-key and secret-key are required (or set KIOOM_APP_KEY/KIOOM_SECRET_KEY)")
	}
	return nil
}

func runAuth(client *kioom.Client, output, action string, _ []string, stdout io.Writer) error {
	switch action {
	case "issue":
		res, err := client.IssueToken()
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "revoke":
		res, err := client.RevokeToken()
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	default:
		return fmt.Errorf("unknown auth action %q", action)
	}
}

func runAccount(client *kioom.Client, output, action string, args []string, stdout io.Writer) error {
	switch action {
	case "number":
		res, err := client.GetAccountNumber()
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "deposit":
		var req kioom.DepositRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateDepositRequest(&req); err != nil {
			return err
		}
		res, err := client.GetDeposit(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "balance":
		var req kioom.AccountBalanceRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateAccountBalanceRequest(&req); err != nil {
			return err
		}
		res, err := client.GetAccountBalance(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	default:
		return fmt.Errorf("unknown account action %q", action)
	}
}

func runStock(client *kioom.Client, output, action string, args []string, stdout io.Writer) error {
	switch action {
	case "basic":
		var req kioom.StockBasicInfoRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateStockCode(req.StkCd); err != nil {
			return err
		}
		res, err := client.GetStockBasicInfo(req.StkCd)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "rank":
		var req kioom.RealtimeItemRankRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateRankRequest(&req); err != nil {
			return err
		}
		res, err := client.GetRealtimeStockRank(req.QryTp)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "minute-chart":
		var req kioom.StockMinuteChartRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateMinuteChartRequest(&req); err != nil {
			return err
		}
		res, err := client.GetStockMinuteChart(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	default:
		return fmt.Errorf("unknown stock action %q", action)
	}
}

func runOrder(client *kioom.Client, output, action string, args []string, stdout io.Writer) error {
	switch action {
	case "buy":
		var req kioom.OrderRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateOrderRequest(&req); err != nil {
			return err
		}
		res, err := client.OrderBuy(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "sell":
		var req kioom.OrderRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateOrderRequest(&req); err != nil {
			return err
		}
		res, err := client.OrderSell(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "modify":
		var req kioom.OrderModifyRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateOrderModifyRequest(&req); err != nil {
			return err
		}
		res, err := client.OrderModify(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	case "cancel":
		var req kioom.OrderCancelRequest
		if err := parseJSONArg(args, &req); err != nil {
			return err
		}
		if err := validateOrderCancelRequest(&req); err != nil {
			return err
		}
		res, err := client.OrderCancel(&req)
		if err != nil {
			return err
		}
		return writeOK(stdout, output, res)
	default:
		return fmt.Errorf("unknown order action %q", action)
	}
}
