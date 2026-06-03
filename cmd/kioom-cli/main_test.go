package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseGlobal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		env       map[string]string
		wantErr   bool
		wantRest  []string
		wantMock  bool
		wantOut   string
		wantApp   string
		wantToken string
	}{
		{
			name:      "reads env by default",
			args:      []string{"schema"},
			env:       map[string]string{"KIOOM_APP_KEY": "ak", "KIOOM_SECRET_KEY": "sk", "KIOOM_TOKEN": "tok"},
			wantRest:  []string{"schema"},
			wantOut:   "json",
			wantApp:   "ak",
			wantToken: "tok",
		},
		{
			name:      "cli flags override env",
			args:      []string{"--app-key", "cli-ak", "--output", "pretty", "--mock", "schema"},
			env:       map[string]string{"KIOOM_APP_KEY": "env-ak"},
			wantRest:  []string{"schema"},
			wantMock:  true,
			wantOut:   "pretty",
			wantApp:   "cli-ak",
			wantToken: "",
		},
		{
			name:    "rejects invalid output",
			args:    []string{"--output", "yaml", "schema"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			getenv := func(k string) string {
				return tt.env[k]
			}
			gotCfg, gotRest, err := parseGlobal(tt.args, getenv)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseGlobal() error = %v", err)
			}
			if strings.Join(gotRest, ",") != strings.Join(tt.wantRest, ",") {
				t.Fatalf("rest mismatch: got=%v want=%v", gotRest, tt.wantRest)
			}
			if gotCfg.mock != tt.wantMock {
				t.Fatalf("mock mismatch: got=%v want=%v", gotCfg.mock, tt.wantMock)
			}
			if gotCfg.output != tt.wantOut {
				t.Fatalf("output mismatch: got=%q want=%q", gotCfg.output, tt.wantOut)
			}
			if gotCfg.appKey != tt.wantApp {
				t.Fatalf("appKey mismatch: got=%q want=%q", gotCfg.appKey, tt.wantApp)
			}
			if gotCfg.token != tt.wantToken {
				t.Fatalf("token mismatch: got=%q want=%q", gotCfg.token, tt.wantToken)
			}
		})
	}
}

func TestRunSchema(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"stock.basic":         "stk_cd",
		"stock.indicators":    "stk_cd",
		"stock.rank":          "qry_tp",
		"stock.volume-surge":  "mrkt_tp",
		"stock.tick-chart":    "stk_cd",
		"stock.minute-chart":  "stk_cd",
		"stock.daily-chart":   "stk_cd",
		"stock.weekly-chart":  "stk_cd",
		"stock.monthly-chart": "stk_cd",
		"stock.yearly-chart":  "stk_cd",
	}
	for tc, field := range cases {
		tc := tc
		field := field
		t.Run(tc, func(t *testing.T) {
			out := new(bytes.Buffer)
			errOut := new(bytes.Buffer)

			code := run([]string{"schema", tc}, out, errOut, func(string) string { return "" })
			if code != 0 {
				t.Fatalf("run() exit code = %d, stderr=%s", code, errOut.String())
			}

			body := out.String()
			if !strings.Contains(body, `"ok":true`) {
				t.Fatalf("expected ok envelope, got %s", body)
			}
			if !strings.Contains(body, `"command":"`+tc+`"`) {
				t.Fatalf("expected command path in body, got %s", body)
			}
			if !strings.Contains(body, `"`+field+`"`) {
				t.Fatalf("expected schema field %q in body, got %s", field, body)
			}
		})
	}
}

func TestRunRequiresCredentialsForApiCalls(t *testing.T) {
	t.Parallel()

	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	code := run([]string{"auth", "issue"}, out, errOut, func(string) string { return "" })
	if code != 2 {
		t.Fatalf("run() exit code = %d, want 2", code)
	}
	if !strings.Contains(errOut.String(), "missing_credentials") {
		t.Fatalf("expected missing_credentials in stderr, got %s", errOut.String())
	}
}

func TestParseJSONArg(t *testing.T) {
	t.Parallel()

	type req struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "valid", args: []string{"--json", `{"name":"abc"}`}},
		{name: "unknown field rejected", args: []string{"--json", `{"name":"abc","x":1}`}, wantErr: true},
		{name: "missing json", args: []string{}, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got req
			err := parseJSONArg(tt.args, &got)
			if tt.wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
