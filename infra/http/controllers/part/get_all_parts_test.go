package part

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParsePaginationBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		query      string
		wantLimit  int64
		wantOffset int64
		wantStatus int
	}{
		{name: "defaults", query: "", wantLimit: 50, wantOffset: 0},
		{name: "explicit values", query: "?limit=10&offset=20", wantLimit: 10, wantOffset: 20},
		{name: "max limit", query: "?limit=200", wantLimit: 200, wantOffset: 0},
		{name: "limit above max", query: "?limit=201", wantStatus: 400},
		{name: "zero limit", query: "?limit=0", wantStatus: 400},
		{name: "negative limit", query: "?limit=-1", wantStatus: 400},
		{name: "non numeric limit", query: "?limit=abc", wantStatus: 400},
		{name: "negative offset", query: "?offset=-1", wantStatus: 400},
		{name: "non numeric offset", query: "?offset=abc", wantStatus: 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest("GET", "/api/parts/"+tt.query, nil)

			args, err := parsePagination(ctx)

			if tt.wantStatus != 0 {
				if err == nil || err.StatusCode != tt.wantStatus {
					t.Fatalf("expected status %d, got %#v", tt.wantStatus, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if args.Limit != tt.wantLimit || args.Offset != tt.wantOffset {
				t.Fatalf("expected limit %d offset %d, got %d/%d", tt.wantLimit, tt.wantOffset, args.Limit, args.Offset)
			}
		})
	}
}
