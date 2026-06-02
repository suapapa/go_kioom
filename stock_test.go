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

func TestGetStockCharts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiID := r.Header.Get("api-id")
		w.Header().Set("Content-Type", "application/json")

		switch apiID {
		case API_GetDailyChart:
			res := StockDailyChartResponse{
				StkCd: "005930",
				StkDtPoleChartQry: []StockChartPole{
					{
						CurPrc:     "70100",
						TrdeQty:    "9263135",
						TrdePrica:  "648525",
						Dt:         "20250908",
						OpenPric:   "69800",
						HighPric:   "70500",
						LowPric:    "69600",
						PredPre:    "+600",
						PredPreSig: "2",
						TrdeTernRt: "+0.16",
					},
				},
				ReturnCode: 0,
			}
			json.NewEncoder(w).Encode(res)
		case API_GetWeeklyChart:
			res := StockWeeklyChartResponse{
				StkCd: "005930",
				StkStkPoleChartQry: []StockChartPole{
					{
						CurPrc:     "69500",
						TrdeQty:    "56700518",
						TrdePrica:  "3922030535087",
						Dt:         "20250901",
						OpenPric:   "68400",
						HighPric:   "70400",
						LowPric:    "67500",
						PredPre:    "-200",
						PredPreSig: "5",
						TrdeTernRt: "+0.96",
					},
				},
				ReturnCode: 0,
			}
			json.NewEncoder(w).Encode(res)
		case API_GetMonthlyChart:
			res := StockMonthlyChartResponse{
				StkCd: "005930",
				StkMthPoleChartQry: []StockChartPole{
					{
						CurPrc:     "78900",
						TrdeQty:    "215040968",
						TrdePrica:  "15774571011618",
						Dt:         "20250901",
						OpenPric:   "68400",
						HighPric:   "79500",
						LowPric:    "67500",
						PredPre:    "+9200",
						PredPreSig: "2",
						TrdeTernRt: "+3.38",
					},
				},
				ReturnCode: 0,
			}
			json.NewEncoder(w).Encode(res)
		case API_GetYearlyChart:
			res := StockYearlyChartResponse{
				StkCd: "005930",
				StkYrPoleChartQry: []StockChartPole{
					{
						CurPrc:    "78200",
						TrdeQty:   "10541142553",
						TrdePrica: "698972287992549",
						Dt:        "20250102",
						OpenPric:  "65100",
						HighPric:  "118800",
						LowPric:   "34900",
					},
				},
				ReturnCode: 0,
			}
			json.NewEncoder(w).Encode(res)
		case API_GetTickChart:
			res := StockTickChartResponse{
				StkCd: "005930",
				StkTicChartQry: []StockTickPole{
					{
						CurPrc:     "78900",
						TrdeQty:    "143",
						CntrTm:     "20250917131939",
						OpenPric:   "78900",
						HighPric:   "78900",
						LowPric:    "78900",
						PredPre:    "500",
						PredPreSig: "5",
					},
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

	req := &StockChartRequest{
		StkCd:      "005930",
		BaseDt:     "20250908",
		UpdStkpcTp: "1",
	}

	t.Run("Daily", func(t *testing.T) {
		res, err := client.GetStockDailyChart(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StkCd != "005930" || len(res.StkDtPoleChartQry) != 1 || res.StkDtPoleChartQry[0].CurPrc != "70100" {
			t.Errorf("invalid daily chart response: %+v", res)
		}
	})

	t.Run("Weekly", func(t *testing.T) {
		res, err := client.GetStockWeeklyChart(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StkCd != "005930" || len(res.StkStkPoleChartQry) != 1 || res.StkStkPoleChartQry[0].CurPrc != "69500" {
			t.Errorf("invalid weekly chart response: %+v", res)
		}
	})

	t.Run("Monthly", func(t *testing.T) {
		res, err := client.GetStockMonthlyChart(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StkCd != "005930" || len(res.StkMthPoleChartQry) != 1 || res.StkMthPoleChartQry[0].CurPrc != "78900" {
			t.Errorf("invalid monthly chart response: %+v", res)
		}
	})

	t.Run("Yearly", func(t *testing.T) {
		res, err := client.GetStockYearlyChart(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StkCd != "005930" || len(res.StkYrPoleChartQry) != 1 || res.StkYrPoleChartQry[0].CurPrc != "78200" {
			t.Errorf("invalid yearly chart response: %+v", res)
		}
	})

	t.Run("Tick", func(t *testing.T) {
		tickReq := &StockTickChartRequest{
			StkCd:      "005930",
			TicScope:   "1",
			UpdStkpcTp: "1",
		}
		res, err := client.GetStockTickChart(context.Background(), tickReq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StkCd != "005930" || len(res.StkTicChartQry) != 1 || res.StkTicChartQry[0].CurPrc != "78900" {
			t.Errorf("invalid tick chart response: %+v", res)
		}
	})
}
