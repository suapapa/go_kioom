# 🌊 Project "Vibe Coding" Guide (AGENTS.md)

This document serves as a guide for both human developers and AI agents to maintain the **"Vibe"** of the `go_kioom` project. "Vibe Coding" in this context means producing code that is not just functional, but also idiomatic, elegant, and highly professional in its implementation and developer experience.

## ✨ The Project Philosophy

1.  **Idiomatic Golang**: We prioritize clean, standard, and idiomatic Go. Avoid over-engineering. Favor clarity over cleverness.
2.  **Developer Experience (DX)**: Every API, function, and example should be easy to understand and use at first glance.
3.  **Premium Execution**: No half-baked implementations. Every feature comes with:
    *   Comprehensive unit tests (table-driven where applicable).
    *   Clear doc comments for all exported types and functions.
    *   Practical, runnable examples in the `cmd/examples/` directory.
4.  **Consistency**: Follow the established patterns for HTTP requests, error handling, and project structure.

---

## 🏗 Core Design Patterns

### 1. Adding a New API Endpoint
When adding a new Kiwoom REST API endpoint (e.g., `ka20001` - Order), follow this workflow:

1.  **Define Structures**: Create request and response structs in a new or existing file (e.g., `order.go`). Use `json` tags correctly as per the Kiwoom spec.
    ```go
    type OrderRequest struct {
        // ... fields ...
    }
    type OrderResponse struct {
        // ... fields ...
    }
    ```
2.  **Implement Method**: Add a method to the `Client` struct in the corresponding file. Use `c.newRequest` and `c.do`.
    ```go
    func (c *Client) SendOrder(req *OrderRequest) (*OrderResponse, error) {
        httpReq, err := c.newRequest("POST", "/v1/order", "ka20001", req)
        if err != nil {
            return nil, err
        }
        var res OrderResponse
        if err := c.do(httpReq, &res); err != nil {
            return nil, err
        }
        return &res, nil
    }
    ```
3.  **Authentication Handshake**: Ensure methods that require a token check `c.Token` or trigger an appropriate error.
4.  **Extend programs under `cmd/`**: Whenever you add or materially change functionality in the root library (`/`), ship the same capability in the relevant `cmd/*` binaries so CLI users and agents stay aligned with the SDK. Typical follow-ups:
    * **Agent-facing CLI** (`cmd/kioom-cli/`): wire `<section> <action>`, add `--json` parsing and local validation consistent with existing commands, register the pair in `commandSchemaMap` (see `schema.go`), add or extend table-driven tests (see `main_test.go`), update `README.md` under that command.
    * **Examples**: add or update runnable examples under `cmd/examples/` when the feature is intended for tutorials or copy-paste starting points.
    * **Other `cmd/` tools**: if another entry under `cmd/` (for example MCP) exposes API surface, mirror the new endpoint or behavior there in the same change or a coordinated follow-up so nothing drifts.

### 2. Testing Vibe
We use `http.HandlerFunc` to mock the Kiwoom API in tests. Maintain this pattern to avoid external dependencies during tests.
*   See `client_test.go` or `auth_test.go` for reference.
*   Always use table-driven tests for multiple scenarios (success, auth error, validation error, etc.).

### 3. Keep `cmd/` in sync (general rule)
Any new user-facing capability in the library should land with corresponding updates under `cmd/` wherever that capability is meant to be used (CLI, examples, MCP, and so on). Do not ship library-only changes and leave shipped `cmd/` programs behind the core package.

---

## 🛠 Tooling and Workflow

-   **Go Modules**: Use `go mod tidy` after adding new dependencies.
-   **Formatting**: Always run `go fmt ./...` or use an editor that does it automatically.
-   **Testing**: Run `go test ./...` frequently to ensure no regressions.
-   **Documentation**: Use `godoc` style comments.
    ```go
    // NewClient returns a new Kiwoom REST API client.
    // appKey and secretKey are provided by Kiwoom portal.
    // if useMock is true, it uses the mock investment environment domain.
    func NewClient(appKey, secretKey string, useMock bool) *Client { ... }
    ```

---

## 🎨 Aesthetic Standards for AI Agents

When acting as an agent on this project, ensure:
1.  **Rich Context**: Before implementing, check `_ref/kiwoom-rest-api-spec.json` (official Kiwoom REST API spec), plus any other `_ref/` materials, to ensure field names and types are accurate.
2.  **Self-Correction**: If you notice a pattern mismatch (e.g., inconsistent error naming), proactively fix it to align with the rest of the project.
3.  **WOW Factor**: Don't just implement the minimum. Provide a helpful example or a robust test case that proves the implementation works.
4.  **Documentation Maintenance**: Proactively update `AGENTS.md` and `README.md` when new patterns, features, or significant changes are introduced to keep the documentation in sync with the codebase.

---

## 🗂 Project Structure Breakdown

-   `/`: Core library files (`client.go`, `auth.go`, `account.go`, etc.).
-   `cmd/kioom-cli/`: Agent-oriented CLI; keep parity with library endpoints via `cli.go`, `schema.go`, `validate.go`, and tests. Shares request validation with MCP via `internal/kioomvalidate/`.
-   `cmd/kioom-mcp/`: MCP server entry (`kioom-mcp`); runs the shared server in `internal/mcpkioom/` with `-transport stdio` (default) or `-transport sse`. Add new API tools only in `internal/mcpkioom/` so both transports stay aligned.
-   `internal/mcpkioom/`: Registers Kiwoom tools on a single `mcp.Server` instance consumed by stdio and SSE transports.
-   `internal/kioomvalidate/`: Shared field validations for CLI and MCP tool handlers.
-   `internal/kioomenv/`: Shared `KIOOM_*` environment loading for `kioom-cli`, `kioom-mcp`, and `cmd/examples/*`.
-   `cmd/examples/`: Practical, runnable examples for users.
-   `_ref/`: Reference materials; primary API spec is `kiwoom-rest-api-spec.json` (official Kiwoom REST API export).
-   `.agents/`: Custom AI instructions and skills (e.g., `golang-pro`).

---

**Remember**: The goal is to make `go_kioom` the most reliable and easy-to-use Kiwoom REST API wrapper in the Go ecosystem. Make it shine! 💎
