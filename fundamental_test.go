package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStockForeignInvestor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_GetStockForeignInvestor {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := StockForeignInvestorResponse{
			StkFrgnr: []StockForeignInvestorItem{
				{Dt: "20260603", Wght: "+52.45"},
			},
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	res, err := client.GetStockForeignInvestor(context.Background(), "005930")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.StkFrgnr) != 1 || res.StkFrgnr[0].Wght != "+52.45" {
		t.Errorf("expected Wght +52.45, got %+v", res.StkFrgnr)
	}
}

func TestGetStockIndicators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiID := r.Header.Get("api-id")
		w.Header().Set("Content-Type", "application/json")

		switch apiID {
		case API_GetStockBasicInfo:
			res := StockBasicInfoResponse{
				StkCd:            "005930",
				StkNm:            "Samsung Electronics",
				Mac:              "4000000",
				Per:              "15.42",
				Pbr:              "1.45",
				Roe:              "9.42",
				High250PricPreRt: "-12.50",
				Low250PricPreRt:  "+20.40",
				ForExhRt:         "52.12",
				SaleAmt:          "3000000",
				BusPro:           "300000",
				CupNga:           "240000",
				ReturnCode:       0,
			}
			json.NewEncoder(w).Encode(res)
		case API_GetStockForeignInvestor:
			res := StockForeignInvestorResponse{
				StkFrgnr: []StockForeignInvestorItem{
					{Wght: "+52.45"},
				},
				ReturnCode: 0,
			}
			json.NewEncoder(w).Encode(res)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	res, err := client.GetStockIndicators(context.Background(), "005930")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StkNm != "Samsung Electronics" {
		t.Errorf("expected StkNm 'Samsung Electronics', got %s", res.StkNm)
	}
	if res.ForQotaRt != "+52.45" {
		t.Errorf("expected ForQotaRt '+52.45', got %s", res.ForQotaRt)
	}
	if res.OperatingMargin != "10.00%" {
		t.Errorf("expected OperatingMargin '10.00%%', got %s", res.OperatingMargin)
	}
	if res.NetProfitMargin != "8.00%" {
		t.Errorf("expected NetProfitMargin '8.00%%', got %s", res.NetProfitMargin)
	}
}
