package buildinfoecho

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/ubgo/buildinfo"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	e := echo.New()
	e.GET("/version", Handler())

	srv := httptest.NewServer(e)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/version")
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
	e := echo.New()
	Mount(e)

	srv := httptest.NewServer(e)
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
	e := echo.New()
	Mount(e, WithPath("/api/v1/version"))

	srv := httptest.NewServer(e)
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
	auth := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-Internal-Key") != "secret" {
				return c.NoContent(http.StatusUnauthorized)
			}
			return next(c)
		}
	}

	e := echo.New()
	Mount(e, WithMiddleware(auth))

	srv := httptest.NewServer(e)
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
