package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	kioom "github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
)

func main() {
	mockFlag := flag.Bool("mock", false, "Use mock domain")
	flag.Parse()

	// For example purposes, we use environment variables for AppKey and SecretKey.
	// In a real application, you should handle these securely.
	cfg := kioomenv.Load(os.Getenv)

	if err := cfg.RequireAppKeys(); err != nil {
		log.Println("KIOOM_APP_KEY or KIOOM_SECRET_KEY not set. Using mock mode.")
	}

	ctx := context.Background()

	// Initialize client (mock mode if keys are missing)
	var opts []kioom.Option
	if *mockFlag {
		opts = append(opts, kioom.WithMockDomain())
	}
	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, opts...)

	// Note: In a real scenario, you need to set the Token after authentication.
	// For this example, we assume the token is already set or the server is mocked.

	// 1. Get Account Number
	fmt.Println("--- Account Number ---")
	accRes, err := client.GetAccountNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get account number: %v", err)
	}
	fmt.Printf("Account Number: %s\n", accRes.AcctNo)

	// 2. Get Deposit (예수금)
	fmt.Println("\n--- Deposit (예수금) ---")
	depRes, err := client.GetDeposit(ctx, &kioom.DepositRequest{
		QryTp: "2", // General
	})
	if err != nil {
		log.Fatalf("failed to get deposit: %v", err)
	}
	fmt.Printf("Total Deposit: %s\n", depRes.Deposit)
	fmt.Printf("Orderable Amount: %s\n", depRes.OrderAllowAmt)
	fmt.Printf("Withdrawable Amount: %s\n", depRes.WithdrawableAmt)

	// 3. Get Account Balance (계좌잔고)
	fmt.Println("\n--- Account Balance (계좌잔고/수익률) ---")
	balRes, err := client.GetAccountBalance(ctx, &kioom.AccountBalanceRequest{
		QryTp:     "1",   // Sum
		DmstStkTP: "KRX", // Korea Exchange
	})
	if err != nil {
		log.Fatalf("failed to get account balance: %v", err)
	}
	fmt.Printf("Total Purchase: %s\n", balRes.TotalPurchaseAmt)
	fmt.Printf("Total Eval: %s\n", balRes.TotalEvalAmt)
	fmt.Printf("Total Profit/Loss: %s\n", balRes.TotalEvalPL)
	fmt.Printf("Return Rate: %v%%\n", balRes.TotalProfitRate)

	fmt.Println("\nItems held:")
	for _, item := range balRes.Items {
		fmt.Printf("- [%s] %s: Qty=%s, Current=%s, P/L=%s\n",
			item.StockCode, item.StockName, item.Quantity, item.CurrentPrice, item.EvalPL)
	}
}
