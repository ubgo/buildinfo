package buildinfogin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ubgo/buildinfo"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandler_ReturnsJSON(t *testing.T) {
	r := gin.New()
	r.GET("/version", Handler())

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
	r := gin.New()
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
	r := gin.New()
	Mount(r, WithPath("/api/v1/version"))

	srv := httptest.NewServer(r)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/api/v1/version")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("override path: got %d, want 200", resp.StatusCode)
	}

	resp404, _ := http.Get(srv.URL + "/version")
	defer resp404.Body.Close()
	if resp404.StatusCode != http.StatusNotFound {
		t.Errorf("default path should be 404, got %d", resp404.StatusCode)
	}
}

func TestMount_WithMiddleware_BlocksUnauthorized(t *testing.T) {
	auth := func(c *gin.Context) {
		if c.GetHeader("X-Internal-Key") != "secret" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}

	r := gin.New()
	Mount(r, WithMiddleware(auth))

	srv := httptest.NewServer(r)
	defer srv.Close()

	respNoKey, _ := http.Get(srv.URL + "/version")
	defer respNoKey.Body.Close()
	if respNoKey.StatusCode != http.StatusUnauthorized {
		t.Errorf("missing key: got %d, want 401", respNoKey.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/version", nil)
	req.Header.Set("X-Internal-Key", "secret")
	respWithKey, _ := http.DefaultClient.Do(req)
	defer respWithKey.Body.Close()
	if respWithKey.StatusCode != http.StatusOK {
		t.Errorf("with key: got %d, want 200", respWithKey.StatusCode)
	}
}
