package main

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSessionLogRegistryResolve(t *testing.T) {
	reg := &sessionLogRegistry{}

	t.Run("nil session", func(t *testing.T) {
		if got := reg.resolve(nil); got != "" {
			t.Fatalf("resolve(nil) = %q, want empty", got)
		}
	})

	t.Run("stable id for same server session", func(t *testing.T) {
		ss := &mcp.ServerSession{}
		first := reg.resolve(ss)
		second := reg.resolve(ss)
		if first == "" {
			t.Fatal("expected non-empty generated session id")
		}
		if first != second {
			t.Fatalf("resolve() = %q and %q, want same id", first, second)
		}
	})

	t.Run("distinct ids for different server sessions", func(t *testing.T) {
		a := reg.resolve(&mcp.ServerSession{})
		b := reg.resolve(&mcp.ServerSession{})
		if a == "" || b == "" {
			t.Fatal("expected non-empty generated session ids")
		}
		if a == b {
			t.Fatalf("distinct sessions got same id %q", a)
		}
	})
}
