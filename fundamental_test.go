package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStockInstitutions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_GetStockInstitutions {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := StockInstitutionsResponse{
			Date:        "20260603",
			ClosePric:   "70000",
			FrgnrQotaRt: "52.45",
			ReturnCode:  0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	res, err := client.GetStockInstitutions(context.Background(), "005930")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.FrgnrQotaRt != "52.45" {
		t.Errorf("expected FrgnrQotaRt 52.45, got %s", res.FrgnrQotaRt)
	}
}

func TestGetStockIndicators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiID := r.Header.Get("api-id")
		w.Header().Set("Content-Type", "application/json")

		if apiID == API_GetStockBasicInfo {
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
				SaleAmt:          "3000000", // 3,000,000 (represented in appropriate units)
				BusPro:           "300000",  // 300,000
				CupNga:           "240000",  // 240,000
				ReturnCode:       0,
			}
			json.NewEncoder(w).Encode(res)
		} else if apiID == API_GetStockInstitutions {
			res := StockInstitutionsResponse{
				FrgnrQotaRt: "52.45",
				ReturnCode:  0,
			}
			json.NewEncoder(w).Encode(res)
		} else {
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
	if res.ForQotaRt != "52.45" {
		t.Errorf("expected ForQotaRt '52.45', got %s", res.ForQotaRt)
	}
	// 영업이익률 = (300000 / 3000000) * 100 = 10.00%
	if res.OperatingMargin != "10.00%" {
		t.Errorf("expected OperatingMargin '10.00%%', got %s", res.OperatingMargin)
	}
	// 순이익률 = (240000 / 3000000) * 100 = 8.00%
	if res.NetProfitMargin != "8.00%" {
		t.Errorf("expected NetProfitMargin '8.00%%', got %s", res.NetProfitMargin)
	}
}
