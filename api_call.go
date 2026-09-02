package kioom

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// CallGeneratedAPI invokes a generated stock/ETF/US API by TR code.
// req and res must be pointers to the generated request/response types (for example *Ka10002Request).
func (c *Client) CallGeneratedAPI(ctx context.Context, apiID string, req, res any) error {
	meta, ok := GeneratedAPIRegistry[apiID]
	if !ok {
		return fmt.Errorf("kioom: unknown generated API %q", apiID)
	}
	if req == nil {
		return fmt.Errorf("kioom: request value is nil for %s", apiID)
	}
	if res == nil {
		return fmt.Errorf("kioom: response value is nil for %s", apiID)
	}
	reqVal := reflect.ValueOf(req)
	resVal := reflect.ValueOf(res)
	if reqVal.Kind() != reflect.Ptr || reqVal.IsNil() {
		return fmt.Errorf("kioom: request must be a non-nil pointer for %s", apiID)
	}
	if resVal.Kind() != reflect.Ptr || resVal.IsNil() {
		return fmt.Errorf("kioom: response must be a non-nil pointer for %s", apiID)
	}

	httpReq, err := c.newRequest(ctx, meta.Method, meta.Path, apiID, req)
	if err != nil {
		return err
	}
	return c.do(httpReq, res)
}

// NewGeneratedRequest allocates a new request value for the given TR code.
func NewGeneratedRequest(apiID string) (any, error) {
	meta, ok := GeneratedAPIRegistry[apiID]
	if !ok {
		return nil, fmt.Errorf("kioom: unknown generated API %q", apiID)
	}
	reqType, err := lookupGeneratedType(meta.RequestType)
	if err != nil {
		return nil, err
	}
	return reflect.New(reqType).Interface(), nil
}

// NewGeneratedResponse allocates a new response value for the given TR code.
func NewGeneratedResponse(apiID string) (any, error) {
	meta, ok := GeneratedAPIRegistry[apiID]
	if !ok {
		return nil, fmt.Errorf("kioom: unknown generated API %q", apiID)
	}
	respType, err := lookupGeneratedType(meta.ResponseType)
	if err != nil {
		return nil, err
	}
	return reflect.New(respType).Interface(), nil
}

// CallGeneratedAPIJSON invokes a generated API using JSON request/response bodies.
func (c *Client) CallGeneratedAPIJSON(ctx context.Context, apiID string, reqJSON []byte) ([]byte, error) {
	req, err := NewGeneratedRequest(apiID)
	if err != nil {
		return nil, err
	}
	if len(reqJSON) > 0 {
		if err := json.Unmarshal(reqJSON, req); err != nil {
			return nil, fmt.Errorf("kioom: invalid request JSON for %s: %w", apiID, err)
		}
	}
	res, err := NewGeneratedResponse(apiID)
	if err != nil {
		return nil, err
	}
	if err := c.CallGeneratedAPI(ctx, apiID, req, res); err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func lookupGeneratedType(name string) (reflect.Type, error) {
	t := generatedTypeByName[name]
	if t == nil {
		return nil, fmt.Errorf("kioom: type %q not registered", name)
	}
	return t, nil
}
