package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	kioom "github.com/suapapa/go_kioom"
)

type apiListItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Method       string `json:"method"`
	Description  string `json:"description,omitempty"`
	Path         string `json:"path"`
	RequestType  string `json:"request_type"`
	ResponseType string `json:"response_type"`
}

func cmdAPIList(_ *kioom.Client, output string, args []string, w io.Writer) error {
	filter := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--filter" && i+1 < len(args) {
			filter = strings.ToLower(strings.TrimSpace(args[i+1]))
			i++
		}
	}

	items := make([]apiListItem, 0, len(kioom.GeneratedAPIRegistry))
	for id, meta := range kioom.GeneratedAPIRegistry {
		if filter != "" && !strings.Contains(strings.ToLower(id), filter) && !strings.Contains(strings.ToLower(meta.Name), filter) {
			continue
		}
		items = append(items, apiListItem{
			ID:           id,
			Name:         meta.Name,
			Method:       meta.MethodName,
			Description:  meta.Description,
			Path:         meta.Path,
			RequestType:  meta.RequestType,
			ResponseType: meta.ResponseType,
		})
	}

	return writeOK(w, output, map[string]any{
		"count": len(items),
		"apis":  items,
	})
}

func cmdAPICall(c *kioom.Client, output string, args []string, w io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kioom api call <api-id> [--json '{...}']")
	}
	apiID := strings.ToLower(strings.TrimSpace(args[0]))
	if _, ok := kioom.GeneratedAPIRegistry[apiID]; !ok {
		return fmt.Errorf("unknown generated API %q (use `kioom api list`)", apiID)
	}

	var reqJSON []byte
	for i := 1; i < len(args); i++ {
		if args[i] == "--json" && i+1 < len(args) {
			reqJSON = []byte(args[i+1])
			break
		}
	}

	raw, err := c.CallGeneratedAPIJSON(context.Background(), apiID, reqJSON)
	if err != nil {
		return err
	}
	if output == "pretty" {
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		return writeOK(w, output, v)
	}
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}
