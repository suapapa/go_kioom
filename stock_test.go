package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStockBasicInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := StockBasicInfoResponse{
			StkCd:      "005930",
			StkNm:      "Samsung Electronics",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	res, err := client.GetStockBasicInfo(context.Background(), "005930")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StkCd != "005930" {
		t.Errorf("expected StkCd 005930, got %s", res.StkCd)
	}
	if res.StkNm != "Samsung Electronics" {
		t.Errorf("expected StkNm Samsung Electronics, got %s", res.StkNm)
	}
}

func TestGetRealtimeStockRank(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := RealtimeItemRankResponse{
			ReturnCode: 0,
			ItemInqRank: []RealtimeItemInqRank{
				{StkCd: "005930", StkNm: "Samsung"},
			},
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	res, err := client.GetRealtimeStockRank(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.ItemInqRank) != 1 {
		t.Fatalf("expected 1 item in rank, got %d", len(res.ItemInqRank))
	}
	if res.ItemInqRank[0].StkCd != "005930" {
		t.Errorf("expected StkCd 005930, got %s", res.ItemInqRank[0].StkCd)
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"", 0},
		{"  ", 0},
		{"+123.45", 123.45},
		{"-54.32", -54.32},
		{"1,234,567.89", 1234567.89},
		{"+1,000", 1000},
		{"invalid", 0},
	}

	for _, tc := range tests {
		got := parseFloat(tc.input)
		if got != tc.expected {
			t.Errorf("parseFloat(%q) = %f; expected %f", tc.input, got, tc.expected)
		}
	}
}

func TestGetVolumeSurge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := VolumeSurgeResponse{
			ReturnCode: 0,
			TrdeQtySdnin: []VolumeSurgeItem{
				{StkCd: "005930", StkNm: "Samsung"},
			},
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	req := &VolumeSurgeRequest{
		MrktTp:    "000",
		SortTp:    "1",
		TmTp:      "2",
		TrdeQtyTp: "5",
		Tm:        "",
		StkCnd:    "0",
		PricTp:    "0",
		StexTp:    "3",
	}
	res, err := client.GetVolumeSurge(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.TrdeQtySdnin) != 1 {
		t.Fatalf("expected 1 item, got %d", len(res.TrdeQtySdnin))
	}
	if res.TrdeQtySdnin[0].StkCd != "005930" {
		t.Errorf("expected StkCd 005930, got %s", res.TrdeQtySdnin[0].StkCd)
	}
}
