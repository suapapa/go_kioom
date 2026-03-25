package kioom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := AccountResponse{
			AcctNo:     "1234567890",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL

	res, err := client.GetAccountNumber()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.AcctNo != "1234567890" {
		t.Errorf("expected AcctNo 1234567890, got %s", res.AcctNo)
	}
}

func TestGetDeposit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := DepositResponse{
			Deposit:         "1000000",
			OrderAllowAmt:   "800000",
			WithdrawableAmt: "900000",
			ReturnCode:      0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL

	res, err := client.GetDeposit(&DepositRequest{QryTp: "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Deposit != "1000000" {
		t.Errorf("expected Deposit 1000000, got %s", res.Deposit)
	}
}

func TestGetAccountBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := AccountBalanceResponse{
			AssetAmt: "10000000",
			Items: []AccountBalanceItem{
				{
					StockCode: "005930",
					StockName: "삼성전자",
					Quantity:  "10",
				},
			},
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL

	res, err := client.GetAccountBalance(&AccountBalanceRequest{QryTp: "1", DmstStkTP: "KRX"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.AssetAmt != "10000000" {
		t.Errorf("expected AssetAmt 10000000, got %s", res.AssetAmt)
	}
	if len(res.Items) != 1 || res.Items[0].StockCode != "005930" {
		t.Errorf("unexpected items: %+v", res.Items)
	}
}
