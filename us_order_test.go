package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUSOrderBuy(t *testing.T) {
	expectedOrdNo := "000000282"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USOrderBuy {
			t.Errorf("expected api-id %s, got %s", API_USOrderBuy, r.Header.Get("api-id"))
		}
		if r.URL.Path != "/api/us/ordr" {
			t.Errorf("expected path /api/us/ordr, got %s", r.URL.Path)
		}
		var body USOrderBuyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.StkCd != "NVDA" || body.StexTp != "ND" {
			t.Errorf("unexpected body: %+v", body)
		}
		res := USOrderBuyResponse{
			StkNm:      "엔비디아",
			OrdNo:      expectedOrdNo,
			ReturnCode: 0,
			ReturnMsg:  "미국 매수주문입력 완료되었습니다",
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	req := &USOrderBuyRequest{
		StexTp: "ND",
		StkCd:  "NVDA",
		OrdQty: "10",
		OrdUv:  "213.04",
		TrdeTp: "00",
	}

	res, err := client.USOrderBuy(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
}

func TestUSOrderSell(t *testing.T) {
	expectedOrdNo := "000000283"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USOrderSell {
			t.Errorf("expected api-id %s, got %s", API_USOrderSell, r.Header.Get("api-id"))
		}
		res := USOrderSellResponse{
			StkNm:      "엔비디아",
			OrdNo:      expectedOrdNo,
			PossQty:    "000000000028",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	req := &USOrderSellRequest{
		StkCd:  "NVDA",
		StexTp: "ND",
		OrdQty: "10",
		OrdUv:  "210.05",
		TrdeTp: "00",
	}

	res, err := client.USOrderSell(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
}

func TestUSOrderModify(t *testing.T) {
	expectedOrdNo := "000000284"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USOrderModify {
			t.Errorf("expected api-id %s, got %s", API_USOrderModify, r.Header.Get("api-id"))
		}
		res := USOrderModifyResponse{
			OrdNo:      expectedOrdNo,
			MdfyOrdQty: "000000000001",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	req := &USOrderModifyRequest{
		OrigOrdNo: "000000050",
		StexTp:    "ND",
		StkCd:     "NVDA",
		MdfyUv:    "210",
	}

	res, err := client.USOrderModify(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
}

func TestUSOrderCancel(t *testing.T) {
	expectedOrdNo := "000000285"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USOrderCancel {
			t.Errorf("expected api-id %s, got %s", API_USOrderCancel, r.Header.Get("api-id"))
		}
		res := USOrderCancelResponse{
			OrdNo:      expectedOrdNo,
			CnclOrdQty: "000000000001",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	req := &USOrderCancelRequest{
		OrigOrdNo: "000000047",
		StexTp:    "ND",
		StkCd:     "NVDA",
	}

	res, err := client.USOrderCancel(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OrdNo != expectedOrdNo {
		t.Errorf("expected ord_no %s, got %s", expectedOrdNo, res.OrdNo)
	}
	if res.CnclOrdQty != "000000000001" {
		t.Errorf("expected cncl_ord_qty 000000000001, got %s", res.CnclOrdQty)
	}
}
