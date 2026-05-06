# `kioom` CLI - AI Agent Guide

This document explains how an AI agent should use the `kioom` CLI safely and reliably.

## Purpose

`kioom` is an agent-friendly wrapper around the Go Kiwoom client. It provides:

- deterministic JSON I/O
- strict request validation for known commands
- runtime schema introspection (`schema` command)

## Install

```bash
go install github.com/suapapa/go_kioom/cmd/kioom
```

## Build

```bash
go build -o bin/kioom ./cmd/kioom
```

## Global Options

- `--app-key`: Kiwoom app key (fallback: `KIOOM_APP_KEY`)
- `--secret-key`: Kiwoom secret key (fallback: `KIOOM_SECRET_KEY`)
- `--token`: bearer token (fallback: `KIOOM_TOKEN`)
- `--mock`: use Kiwoom mock endpoint
- `--output json|pretty`: output format (default: `json`)

Credential rules:

- `app-key` and `secret-key` are required for all non-`schema` commands.
- `token` is required for most account/stock/order operations.
- `auth issue` is used to obtain a token.

## Output Contract (Machine-Readable)

All command outputs are JSON envelopes.

Success:

```json
{"ok":true,"data":{...}}
```

Failure:

```json
{"ok":false,"error":{"code":"...","message":"..."}}
```

Agent behavior recommendation:

1. Parse JSON.
2. Check `ok`.
3. If `ok=false`, branch by `error.code` first, then inspect `error.message`.

Common error codes:

- `invalid_arguments`
- `missing_credentials`
- `command_failed`
- `unknown_schema_path`
- `schema_error`
- `output_error`
- `usage`

Exit code conventions:

- `0`: success
- `1`: runtime/command failure
- `2`: invalid usage or validation/argument failure

## Input Contract for Action Commands

Commands that require request payloads accept them through:

- `--json '<request-json>'`

Rules:

- Payload must match request struct fields exactly.
- Unknown JSON fields are rejected.
- Multiple JSON values are rejected.
- Some fields are validated locally (for example stock code format and enum-like fields).

## Command Layout

```bash
kioom [global-options] <section> <action> [options]
kioom [global-options] schema <command-path>
```

Supported sections/actions:

- `auth issue|revoke`
- `account number|deposit|balance`
- `stock basic|rank|minute-chart`
- `order buy|sell|modify|cancel`

## Schema Introspection (Strongly Recommended for Agents)

List available command paths:

```bash
kioom schema
```

Inspect a specific command schema:

```bash
kioom schema order.buy
```

Use `schema` when:

- generating a request body automatically
- validating mappings before execution
- adapting to command-specific request/response structures

## Suggested Agent Execution Flow

1. Ensure binary exists (`bin/kioom`) or build it.
2. Call `kioom schema` and select the target command path.
3. Call `kioom schema <command-path>` and build `--json` payload from request fields.
4. Issue token if needed:
   - `kioom --mock auth issue` (or without `--mock` for production)
5. Execute target command with global credentials and token.
6. Parse envelope JSON, verify `ok`, then consume `data`.

## Practical Examples

```bash
# 1) Issue token (mock environment)
kioom --mock auth issue

# 2) Get account number
kioom --mock --token "$KIOOM_TOKEN" account number

# 3) Get stock basic info
kioom --mock stock basic --json '{"stk_cd":"005930"}'

# 4) Submit buy order
kioom --mock order buy --json '{"dmst_stex_tp":"KRX","stk_cd":"005930","ord_qty":"1","ord_uv":"0","trde_tp":"3","cond_uv":""}'

# 5) Inspect request/response schema
kioom schema stock.minute-chart
```
