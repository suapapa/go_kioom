package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallGeneratedAPIJSON_GetStockBroker(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != "ka10002" {
			t.Fatalf("api-id: got %q", r.Header.Get("api-id"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"stk_cd":"005930","stk_nm":"삼성전자","return_code":0,"return_msg":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("app", "secret")
	client.baseURL = server.URL
	client.SetToken("test-token")

	raw, err := client.CallGeneratedAPIJSON(context.Background(), "ka10002", []byte(`{"stk_cd":"005930"}`))
	if err != nil {
		t.Fatalf("CallGeneratedAPIJSON: %v", err)
	}

	var res StockBrokerResponse
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if res.StkCd != "005930" {
		t.Fatalf("stk_cd: got %q", res.StkCd)
	}
	if res.ReturnCode != 0 {
		t.Fatalf("return_code: got %d", res.ReturnCode)
	}
}

func TestGeneratedAPIRegistryContainsGetStockBroker(t *testing.T) {
	t.Parallel()

	meta, ok := GeneratedAPIRegistry["ka10002"]
	if !ok {
		t.Fatal("ka10002 not in registry")
	}
	if meta.MethodName != "GetStockBroker" {
		t.Fatalf("method name: got %q", meta.MethodName)
	}
	if meta.RequestType != "StockBrokerRequest" {
		t.Fatalf("request type: got %q", meta.RequestType)
	}
}
