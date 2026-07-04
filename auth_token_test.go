package kioom

import (
	"testing"
	"time"
)

func TestParseExpiresDt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "compact timestamp",
			input: "20241107083713",
		},
		{
			name:  "spaced timestamp",
			input: "2024-11-07 08:37:13",
		},
		{
			name:    "empty",
			input:   "",
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseExpiresDt(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.input == "" && !got.IsZero() {
				t.Fatalf("expected zero time for empty input")
			}
			if tt.input != "" && got.IsZero() {
				t.Fatalf("expected parsed time, got zero")
			}
		})
	}
}

func TestIsTokenAuthError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		meta apiResponseMeta
		want bool
	}{
		{
			name: "success",
			meta: apiResponseMeta{ReturnCode: 0, ReturnMsg: "ok"},
			want: false,
		},
		{
			name: "return code 3",
			meta: apiResponseMeta{ReturnCode: 3, ReturnMsg: "8005:Token이 유효하지 않습니다"},
			want: true,
		},
		{
			name: "8005 in message",
			meta: apiResponseMeta{ReturnCode: 99, ReturnMsg: "8005:Token이 유효하지 않습니다"},
			want: true,
		},
		{
			name: "other business error",
			meta: apiResponseMeta{ReturnCode: 1, ReturnMsg: "invalid parameter"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isTokenAuthError(tt.meta); got != tt.want {
				t.Fatalf("isTokenAuthError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenValidLocked(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	client := NewClient("app", "sec")
	client.setTokenWithExpiry("token", now.Add(30*time.Minute))

	if !client.tokenValidAt(now) {
		t.Fatal("expected token to be valid before expiry")
	}

	client.setTokenWithExpiry("token", now.Add(30*time.Second))
	if client.tokenValidAt(now) {
		t.Fatal("expected token inside expiry buffer to be treated as expired")
	}
}
