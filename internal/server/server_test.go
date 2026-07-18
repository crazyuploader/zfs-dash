package server

import (
	"net/http/httptest"
	"testing"

	"github.com/crazyuploader/zfs-dash/internal/model"
	"github.com/gofiber/fiber/v3"
)

func testNodes() []model.NodeData {
	return []model.NodeData{{
		Label: "n1",
		URL:   "http://internal:9134/metrics",
		Error: "boom",
	}}
}

func TestParseHistoryQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
		from    int64
		to      int64
		bucket  int64
	}{
		{"empty", "", false, 0, 0, 0},
		{"valid range", "from=100&to=200&bucket=60", false, 100, 200, 60},
		{"invalid from", "from=abc", true, 0, 0, 0},
		{"invalid to", "to=abc", true, 0, 0, 0},
		{"invalid bucket", "bucket=abc", true, 0, 0, 0},
		{"negative from", "from=-1", true, 0, 0, 0},
		{"negative bucket", "bucket=-60", true, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			var from, to, bucket int64
			var err error
			app.Get("/t", func(c fiber.Ctx) error {
				from, to, bucket, err = parseHistoryQueryParams(c)
				return nil
			})
			req := httptest.NewRequest("GET", "/t?"+tt.query, nil)
			if _, terr := app.Test(req); terr != nil {
				t.Fatalf("app.Test: %v", terr)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil {
				if from != tt.from || to != tt.to || bucket != tt.bucket {
					t.Errorf("got (%d,%d,%d), want (%d,%d,%d)", from, to, bucket, tt.from, tt.to, tt.bucket)
				}
			}
		})
	}
}

func TestNodeViewsStripURL(t *testing.T) {
	views := nodeViews(testNodes())
	if len(views) != 1 {
		t.Fatalf("expected 1 view, got %d", len(views))
	}
	v := views[0]
	if v.Label != "n1" || v.Error != "boom" {
		t.Errorf("unexpected view: %+v", v)
	}
}
