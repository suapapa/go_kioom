package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUSOpenOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USOpenOrders {
			t.Errorf("expected api-id %s, got %s", API_USOpenOrders, r.Header.Get("api-id"))
		}
		if r.URL.Path != "/api/us/acnt" {
			t.Errorf("expected path /api/us/acnt, got %s", r.URL.Path)
		}
		res := USOpenOrdersResponse{
			ResultList: []USOpenOrderItem{
				{
					OrdNo:     "000000282",
					StkCd:     "NVDA",
					FrgnStkNm: "엔비디아",
					SlbyTpNm:  "매수",
					OrdStat:   "접수",
				},
			},
			ReturnCode: 0,
			ReturnMsg:  "계좌별 미체결내역이 조회되었습니다.",
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	res, err := client.GetUSOpenOrders(context.Background(), &USOpenOrdersRequest{SlbyTp: "0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.ResultList) != 1 {
		t.Fatalf("expected 1 item, got %d", len(res.ResultList))
	}
	if res.ResultList[0].StkCd != "NVDA" {
		t.Errorf("expected stk_cd NVDA, got %s", res.ResultList[0].StkCd)
	}
}

func TestGetUSAccountBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USAccountBalance {
			t.Errorf("expected api-id %s, got %s", API_USAccountBalance, r.Header.Get("api-id"))
		}
		res := USAccountBalanceResponse{
			CrncCode:   "USD",
			TotEvltAmt: "156464.6701",
			ResultList: []USAccountBalanceItem{
				{
					StkCd:     "AAPL",
					FrgnStkNm: "애플",
					PossQty:   "000000000395",
				},
			},
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	res, err := client.GetUSAccountBalance(context.Background(), &USAccountBalanceRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotEvltAmt != "156464.6701" {
		t.Errorf("unexpected tot_evlt_amt: %s", res.TotEvltAmt)
	}
	if len(res.ResultList) != 1 || res.ResultList[0].StkCd != "AAPL" {
		t.Errorf("unexpected result_list: %+v", res.ResultList)
	}
}

func TestGetUSDeposit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_USDeposit {
			t.Errorf("expected api-id %s, got %s", API_USDeposit, r.Header.Get("api-id"))
		}
		res := USDepositResponse{
			KrwEntra: "000000930907881",
			ResultList: []USDepositItem{
				{
					CrncCode: "USD",
					CrncNm:   "미국달러",
					FcEntra:  "18039493.57",
				},
			},
			ReturnCode: 0,
			ReturnMsg:  "조회가 완료되었습니다",
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("valid-token")

	res, err := client.GetUSDeposit(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KrwEntra != "000000930907881" {
		t.Errorf("unexpected krw_entra: %s", res.KrwEntra)
	}
	if len(res.ResultList) != 1 || res.ResultList[0].CrncCode != "USD" {
		t.Errorf("unexpected result_list: %+v", res.ResultList)
	}
}
