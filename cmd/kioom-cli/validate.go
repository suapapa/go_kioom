package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"
)

func parseJSONArg(args []string, out interface{}) error {
	fs := flag.NewFlagSet("json", flag.ContinueOnError)
	fs.SetOutput(bytes.NewBuffer(nil))

	jsonStr := ""
	fs.StringVar(&jsonStr, "json", "", "Request body in JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(jsonStr) == "" {
		return errors.New("--json is required")
	}

	dec := json.NewDecoder(strings.NewReader(jsonStr))
	dec.DisallowUnknownFields()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("invalid --json payload: %w", err)
	}
	if dec.More() {
		return errors.New("invalid --json payload: multiple JSON values")
	}
	return nil
}
