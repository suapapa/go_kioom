package mcpkioom

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	kioom "github.com/suapapa/go_kioom"
)

func TestNewServerToolCatalog(t *testing.T) {
	t.Parallel()

	client := kioom.NewClient("app", "sec", kioom.WithMockDomain())
	server := NewServer(client, &mcp.Implementation{Name: "test-kioom", Version: "0.0.1"}, nil)

	ct, st := mcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}

	cli := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	session, err := cli.Connect(context.Background(), ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	res, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if got := len(res.Tools); got != 19 {
		t.Fatalf("tool count: got %d want 19", got)
	}
}

func TestStockBasicInvalidCodeIsToolError(t *testing.T) {
	t.Parallel()

	client := kioom.NewClient("app", "sec", kioom.WithMockDomain())
	server := NewServer(client, &mcp.Implementation{Name: "test-kioom", Version: "0.0.1"}, nil)

	ct, st := mcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}

	cli := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	session, err := cli.Connect(context.Background(), ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "stock_basic",
		Arguments: map[string]any{"stk_cd": "12"},
	})
	if err != nil {
		t.Fatalf("CallTool protocol error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected validation failure as tool error")
	}
}
