package redirect

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/analytics"
)

type stubResolver struct{ err error }

func (s stubResolver) ResolveAlias(*gin.Context, RedirectDto) (string, error) {
	return "https://example.com", s.err
}

type rejectingPublisher struct{ events []analytics.ClickEventDto }

func (p *rejectingPublisher) TryPublish(event analytics.ClickEventDto) bool {
	p.events = append(p.events, event)
	return false
}

func TestRedirectAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name   string
		alias  string
		err    error
		status int
		events int
	}{
		{"full queue", "abcdef", nil, 302, 1},
		{"missing link", "abcdef", errors.New("not found"), 400, 0},
		{"invalid alias", "invalid", nil, 400, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			publisher := &rejectingPublisher{}
			controller := NewRedirectControllerWithAnalytics(stubResolver{tc.err}, publisher)
			router := gin.New()
			router.GET("/:alias", controller.Redirect)
			req := httptest.NewRequest("GET", "/"+tc.alias, nil)
			req.Header.Set("Referer", "https://referrer.example")
			req.Header.Set("User-Agent", "test-agent")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tc.status || len(publisher.events) != tc.events {
				t.Fatalf("status=%d events=%d", response.Code, len(publisher.events))
			}
			if tc.events > 0 {
				event := publisher.events[0]
				if event.Alias != tc.alias || event.Timestamp.IsZero() || event.Referrer != req.Referer() || event.UserAgent != req.UserAgent() {
					t.Fatalf("unexpected event: %+v", event)
				}
				if response.Header().Get("Location") != "https://example.com" {
					t.Fatal("wrong destination")
				}
			}
		})
	}
}
