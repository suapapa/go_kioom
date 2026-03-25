// Package kioom provides a Go wrapper for the Kiwoom Open API (REST).
//
// This package allows developers to interact with Kiwoom's trading platform
// using idiomatic Go. It handles the complexities of authentication (OAuth2),
// request signing, and data parsing.
//
// Basic usage:
//
//	client := kioom.NewClient("your_app_key", "your_secret_key", true) // true for mock
//	_, err := client.IssueToken()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	res, err := client.GetAccountNumber()
//	if err == nil {
//	    fmt.Println("Account Number:", res.AcctNo)
//	}
//
// For more detailed examples, see the cmd/examples directory.
package kioom
