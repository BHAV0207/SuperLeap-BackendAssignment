package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/observability"
)

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID(), Metrics())
	r.GET("/leads/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestRequestIDGeneratedAndEchoed(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/leads/1", nil))
	if w.Header().Get(RequestIDHeader) == "" {
		t.Fatal("expected a generated request id")
	}

	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/leads/1", nil)
	req.Header.Set(RequestIDHeader, "abc-123")
	r.ServeHTTP(w, req)
	if got := w.Header().Get(RequestIDHeader); got != "abc-123" {
		t.Fatalf("request id = %q, want abc-123", got)
	}
}

func TestMetricsUseRouteTemplate(t *testing.T) {
	r := newRouter()

	counter := observability.HTTPRequestsTotal.WithLabelValues("GET", "/leads/:id", "200")
	before := testutil.ToFloat64(counter)

	for _, id := range []string{"a", "b", "c"} {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/leads/"+id, nil))
	}

	if got := testutil.ToFloat64(counter) - before; got != 3 {
		t.Fatalf("counter increased by %v, want 3", got)
	}

	unmatched := observability.HTTPRequestsTotal.WithLabelValues("GET", "unmatched", "404")
	before = testutil.ToFloat64(unmatched)
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nope", nil))
	if got := testutil.ToFloat64(unmatched) - before; got != 1 {
		t.Fatalf("unmatched counter increased by %v, want 1", got)
	}
}
