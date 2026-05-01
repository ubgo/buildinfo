package buildinfochi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/ubgo/buildinfo"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var info buildinfo.Info
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, body)
	}
	if info.GoVersion == "" {
		t.Errorf("GoVersion empty")
	}
}

func TestMount_DefaultPath(t *testing.T) {
	r := chi.NewRouter()
	Mount(r)

	srv := httptest.NewServer(r)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", resp.StatusCode)
	}
}

func TestMount_WithPathOverride(t *testing.T) {
	r := chi.NewRouter()
	Mount(r, WithPath("/api/v1/version"))

	srv := httptest.NewServer(r)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/version")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("override path: got %d, want 200", resp.StatusCode)
	}

	resp404, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp404.Body.Close()
	if resp404.StatusCode != http.StatusNotFound {
		t.Errorf("default path should be 404, got %d", resp404.StatusCode)
	}
}

func TestMount_WithMiddleware_BlocksUnauthorized(t *testing.T) {
	auth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Internal-Key") != "secret" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	r := chi.NewRouter()
	Mount(r, WithMiddleware(auth))

	srv := httptest.NewServer(r)
	defer srv.Close()

	respNoKey, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer respNoKey.Body.Close()
	if respNoKey.StatusCode != http.StatusUnauthorized {
		t.Errorf("missing key: got %d, want 401", respNoKey.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/version", nil)
	req.Header.Set("X-Internal-Key", "secret")
	respWithKey, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer respWithKey.Body.Close()
	if respWithKey.StatusCode != http.StatusOK {
		t.Errorf("with key: got %d, want 200", respWithKey.StatusCode)
	}
}
