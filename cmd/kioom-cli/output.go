package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type okEnvelope struct {
	OK   bool `json:"ok"`
	Data any  `json:"data"`
}

type errEnvelope struct {
	OK    bool `json:"ok"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeOK(w io.Writer, output string, data any) error {
	return writeEnvelope(w, output, okEnvelope{
		OK:   true,
		Data: data,
	})
}

func writeErr(w io.Writer, output, code string, err error) {
	e := errEnvelope{OK: false}
	e.Error.Code = code
	e.Error.Message = err.Error()
	_ = writeEnvelope(w, output, e)
}

func writeEnvelope(w io.Writer, output string, payload any) error {
	var (
		b   []byte
		err error
	)
	if output == "pretty" {
		b, err = json.MarshalIndent(payload, "", "  ")
	} else {
		b, err = json.Marshal(payload)
	}
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	if _, err := w.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
