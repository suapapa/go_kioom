package kioom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderBuy(t *testing.T) {
	expectedOrdNo := "0001234"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != "kt10000" {
			t.Errorf("expected api-id kt10000, got %s", r.Header.Get("api-id"))
		}
		res := OrderResponse{
			OrdNo:      expectedOrdNo,
			DmstStexTp: "KRX",
			ReturnCode: 0,
			ReturnMsg:  "Order success",
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL
	client.SetToken("valid-token")

	req := &OrderRequest{
		DmstStexTp: "KRX",
		StkCd:      "005930",
		OrdQty:     "10",
		OrdUv:      "70000",
		TrdeTp:     "0",
	}

	res, err := client.OrderBuy(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
}

func TestOrderCancel(t *testing.T) {
	expectedOrdNo := "0001235"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != "kt10003" {
			t.Errorf("expected api-id kt10003, got %s", r.Header.Get("api-id"))
		}
		res := OrderCancelResponse{
			OrdNo:         expectedOrdNo,
			BaseOrigOrdNo: "0001234",
			CnclQty:       "5",
			ReturnCode:    0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL
	client.SetToken("valid-token")

	req := &OrderCancelRequest{
		DmstStexTp: "KRX",
		OrigOrdNo:  "0001234",
		StkCd:      "005930",
		CnclQty:    "5",
	}

	res, err := client.OrderCancel(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
	if res.CnclQty != "5" {
		t.Errorf("expected cncl_qty 5, got %s", res.CnclQty)
	}
}

func TestCreditOrderSell(t *testing.T) {
	expectedOrdNo := "0001236"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != "kt10007" {
			t.Errorf("expected api-id kt10007, got %s", r.Header.Get("api-id"))
		}
		res := OrderResponse{
			OrdNo:      expectedOrdNo,
			DmstStexTp: "KRX",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL
	client.SetToken("valid-token")

	req := &CreditSellRequest{
		DmstStexTp: "KRX",
		StkCd:      "005930",
		OrdQty:     "10",
		OrdUv:      "70000",
		TrdeTp:     "0",
		CrdDealTp:  "33",
		CrdLoanDt:  "20231027",
	}

	res, err := client.CreditOrderSell(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
}
