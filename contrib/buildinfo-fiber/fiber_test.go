package buildinfofiber

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/ubgo/buildinfo"
)

// httpRequest runs a request through the in-process Fiber app and returns the
// status code and body, hiding Fiber's Test() boilerplate.
func httpRequest(t *testing.T, app *fiber.App, method, path string, headers map[string]string) (int, []byte) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func TestHandler_ReturnsJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/version", Handler())

	status, body := httpRequest(t, app, http.MethodGet, "/version", nil)
	if status != http.StatusOK {
		t.Errorf("status: got %d, want 200", status)
	}
	var info buildinfo.Info
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, body)
	}
	if info.GoVersion == "" {
		t.Errorf("GoVersion empty")
	}
}

func TestMount_DefaultPath(t *testing.T) {
	app := fiber.New()
	Mount(app)

	status, _ := httpRequest(t, app, http.MethodGet, "/version", nil)
	if status != http.StatusOK {
		t.Errorf("status: got %d, want 200", status)
	}
}

func TestMount_WithPathOverride(t *testing.T) {
	app := fiber.New()
	Mount(app, WithPath("/api/v1/version"))

	statusOK, _ := httpRequest(t, app, http.MethodGet, "/api/v1/version", nil)
	if statusOK != http.StatusOK {
		t.Errorf("override path: got %d, want 200", statusOK)
	}

	status404, _ := httpRequest(t, app, http.MethodGet, "/version", nil)
	if status404 != http.StatusNotFound {
		t.Errorf("default path should be 404, got %d", status404)
	}
}

func TestMount_WithMiddleware_BlocksUnauthorized(t *testing.T) {
	auth := func(c *fiber.Ctx) error {
		if c.Get("X-Internal-Key") != "secret" {
			return c.SendStatus(http.StatusUnauthorized)
		}
		return c.Next()
	}

	app := fiber.New()
	Mount(app, WithMiddleware(auth))

	statusNoKey, _ := httpRequest(t, app, http.MethodGet, "/version", nil)
	if statusNoKey != http.StatusUnauthorized {
		t.Errorf("missing key: got %d, want 401", statusNoKey)
	}

	statusWithKey, _ := httpRequest(t, app, http.MethodGet, "/version", map[string]string{
		"X-Internal-Key": "secret",
	})
	if statusWithKey != http.StatusOK {
		t.Errorf("with key: got %d, want 200", statusWithKey)
	}
}
