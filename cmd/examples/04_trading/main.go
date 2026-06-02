// Example of sending an order to the Kiwoom REST API.
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

	cfg := kioomenv.Load(os.Getenv)
	if err := cfg.RequireAppKeys(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Use mock domain for testing if needed
	var opts []kioom.Option
	if *mockFlag {
		opts = append(opts, kioom.WithMockDomain())
	}
	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, opts...)

	// 2. Obtain an access token
	tokenRes, err := client.IssueToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}
	fmt.Printf("Successfully authenticated. Token expires at %s.\n", tokenRes.ExpiresDt)
	client.SetToken(tokenRes.Token)

	// 3. Send a Buy Order (Market Price)
	// See: kt10000
	buyReq := &kioom.OrderRequest{
		DmstStexTp: "KRX",
		StkCd:      "005930", // Samsung Electronics
		OrdQty:     "1",
		TrdeTp:     "3", // Market Price
	}

	buyRes, err := client.OrderBuy(ctx, buyReq)
	if err != nil {
		log.Fatalf("Failed to send buy order: %v", err)
	}

	if buyRes.ReturnCode == 0 {
		fmt.Printf("Buy order success! Order Number: %s\n", buyRes.OrdNo)
	} else {
		fmt.Printf("Buy order failed: [%d] %s\n", buyRes.ReturnCode, buyRes.ReturnMsg)
	}

	// 4. Send a Credit Sell Order (Normal Price) if you have any
	// See: kt10007
	// This is just an example, it will likely fail unless you have credit positions.
	creditSellReq := &kioom.CreditSellRequest{
		DmstStexTp: "KRX",
		StkCd:      "005930",
		OrdQty:     "1",
		OrdUv:      "70000",
		TrdeTp:     "0",        // Normal Price
		CrdDealTp:  "33",       // Loan
		CrdLoanDt:  "20231001", // Example loan date
	}

	creditRes, err := client.CreditOrderSell(ctx, creditSellReq)
	if err != nil {
		log.Printf("Credit sell order failed as expected (or due to error): %v\n", err)
	} else {
		fmt.Printf("Credit sell response: [%d] %s\n", creditRes.ReturnCode, creditRes.ReturnMsg)
	}
}
