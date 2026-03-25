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

### 2. Testing Vibe
We use `http.HandlerFunc` to mock the Kiwoom API in tests. Maintain this pattern to avoid external dependencies during tests.
*   See `client_test.go` or `auth_test.go` for reference.
*   Always use table-driven tests for multiple scenarios (success, auth error, validation error, etc.).

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
1.  **Rich Context**: Before implementing, check `kiwoom_api.txt` or related specifications to ensure field names and types are accurate.
2.  **Self-Correction**: If you notice a pattern mismatch (e.g., inconsistent error naming), proactively fix it to align with the rest of the project.
3.  **WOW Factor**: Don't just implement the minimum. Provide a helpful example or a robust test case that proves the implementation works.

---

## 🗂 Project Structure Breakdown

-   `/`: Core library files (`client.go`, `auth.go`, `account.go`, etc.).
-   `cmd/examples/`: Practical, runnable examples for users.
-   `_ref/`: (Optional) Reference documents or original API specs.
-   `.agents/`: Custom AI instructions and skills (e.g., `golang-pro`).

---

**Remember**: The goal is to make `go_kioom` the most reliable and easy-to-use Kiwoom REST API wrapper in the Go ecosystem. Make it shine! 💎
