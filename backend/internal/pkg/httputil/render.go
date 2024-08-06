package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, v interface{}) error {
	defer io.Copy(io.Discard, r.Body) //nolint:errcheck
	return json.NewDecoder(r.Body).Decode(v)
}

// RenderJSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func RenderJSON(w http.ResponseWriter, code int, v interface{}) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("failed to marshal response json: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(buf.Bytes()) //nolint:errcheck
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	w.WriteHeader(code)
	return nil
}
