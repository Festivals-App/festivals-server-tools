package servertools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// --- RespondCode ---

func TestRespondCode(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{"200 OK", http.StatusOK},
		{"201 Created", http.StatusCreated},
		{"400 Bad Request", http.StatusBadRequest},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondCode(w, tt.code)

			if w.Code != tt.code {
				t.Errorf("expected status %d, got %d", tt.code, w.Code)
			}
			if ct := w.Header().Get("Content-Type"); ct != "text/plain" {
				t.Errorf("expected Content-Type text/plain, got %q", ct)
			}
		})
	}
}

// --- RespondString ---

func TestRespondString(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{"simple message", http.StatusOK, "hello world"},
		{"empty message", http.StatusOK, ""},
		{"error message", http.StatusBadRequest, "bad input"},
		{"unicode", http.StatusOK, "héllo wörld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondString(w, tt.code, tt.message)

			if w.Code != tt.code {
				t.Errorf("expected status %d, got %d", tt.code, w.Code)
			}
			if ct := w.Header().Get("Content-Type"); ct != "text/plain" {
				t.Errorf("expected Content-Type text/plain, got %q", ct)
			}
			if body := w.Body.String(); body != tt.message {
				t.Errorf("expected body %q, got %q", tt.message, body)
			}
		})
	}
}

// --- RespondJSON ---

func TestRespondJSON(t *testing.T) {
	t.Run("struct payload", func(t *testing.T) {
		w := httptest.NewRecorder()
		payload := map[string]string{"key": "value"}
		RespondJSON(w, http.StatusOK, payload)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %q", ct)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		data, ok := result["data"]
		if !ok {
			t.Fatal("expected 'data' key in response")
		}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			t.Fatalf("expected data to be a map, got %T", data)
		}
		if dataMap["key"] != "value" {
			t.Errorf("expected data.key = 'value', got %v", dataMap["key"])
		}
	})

	t.Run("empty slice payload becomes []", func(t *testing.T) {
		w := httptest.NewRecorder()
		// fmt.Sprint([]string{}) == "[]"
		RespondJSON(w, http.StatusOK, []string{})

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		data, ok := result["data"]
		if !ok {
			t.Fatal("expected 'data' key")
		}
		arr, ok := data.([]interface{})
		if !ok {
			t.Fatalf("expected data to be []interface{}, got %T", data)
		}
		if len(arr) != 0 {
			t.Errorf("expected empty array, got %v", arr)
		}
	})

	t.Run("non-empty slice payload", func(t *testing.T) {
		w := httptest.NewRecorder()
		RespondJSON(w, http.StatusCreated, []string{"a", "b"})

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		arr := result["data"].([]interface{})
		if len(arr) != 2 {
			t.Errorf("expected 2 items, got %d", len(arr))
		}
	})

	t.Run("nil payload", func(t *testing.T) {
		w := httptest.NewRecorder()
		RespondJSON(w, http.StatusOK, nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if _, ok := result["data"]; !ok {
			t.Error("expected 'data' key in response")
		}
	})
}

// --- RespondError ---

func TestRespondError(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{"404 not found", http.StatusNotFound, "resource not found"},
		{"400 bad request", http.StatusBadRequest, "invalid input"},
		{"500 server error", http.StatusInternalServerError, "something went wrong"},
		{"empty message", http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondError(w, tt.code, tt.message)

			if w.Code != tt.code {
				t.Errorf("expected status %d, got %d", tt.code, w.Code)
			}
			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("expected application/json, got %q", ct)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			errMsg, ok := result["error"]
			if !ok {
				t.Fatal("expected 'error' key in response")
			}
			if errMsg != tt.message {
				t.Errorf("expected error %q, got %q", tt.message, errMsg)
			}
		})
	}
}

// --- RespondFile ---

func TestRespondFile(t *testing.T) {
	t.Run("valid text file", func(t *testing.T) {
		f, err := os.CreateTemp("", "servertools-test-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())

		content := "hello from file"
		if _, err := f.WriteString(content); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}
		// Seek back to start so RespondFile can read it
		if _, err := f.Seek(0, 0); err != nil {
			t.Fatalf("failed to seek: %v", err)
		}

		w := httptest.NewRecorder()
		RespondFile(w, f)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if cl := w.Header().Get("Content-Length"); cl == "" {
			t.Error("expected Content-Length header to be set")
		}
		if ct := w.Header().Get("Content-Type"); ct == "" {
			t.Error("expected Content-Type header to be set")
		}
		if body := w.Body.String(); !strings.Contains(body, content) {
			t.Errorf("expected body to contain %q, got %q", content, body)
		}
	})

	t.Run("closed file returns 500", func(t *testing.T) {
		f, err := os.CreateTemp("", "servertools-test-closed-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		name := f.Name()
		f.Close()
		os.Remove(name)

		// Re-open to get a valid *os.File handle, then close it to simulate error
		// Instead: create a file, close it, and pass the closed handle
		f2, _ := os.CreateTemp("", "servertools-test-closed2-*.txt")
		f2.Close()
		defer os.Remove(f2.Name())

		w := httptest.NewRecorder()
		RespondFile(w, f2)

		// Stat on a closed file will fail -> expect 500
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for closed file, got %d", w.Code)
		}
	})
}

// --- UnauthorizedResponse ---

func TestUnauthorizedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	UnauthorizedResponse(w)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	wwwAuth := w.Header().Get("WWW-Authenticate")
	if !strings.Contains(wwwAuth, "Basic") {
		t.Errorf("expected WWW-Authenticate to contain 'Basic', got %q", wwwAuth)
	}
	if !strings.Contains(wwwAuth, "restricted") {
		t.Errorf("expected WWW-Authenticate to contain 'restricted', got %q", wwwAuth)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if result["error"] != http.StatusText(http.StatusUnauthorized) {
		t.Errorf("expected error %q, got %v", http.StatusText(http.StatusUnauthorized), result["error"])
	}
}
