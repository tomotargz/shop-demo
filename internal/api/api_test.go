package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/shop/internal/catalog"
	"example.com/shop/internal/payment"
)

func TestOrders(t *testing.T) {
	h := newTestHandler()

	t.Run("注文の金額を返す", func(t *testing.T) {
		rec := post(t, h, "/orders", `{"items":[{"sku":"MUG-01","quantity":1}]}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body)
		}
		var got struct{ Total int }
		decode(t, rec, &got)
		if got.Total != 2198 {
			t.Errorf("Total = %d, want 2198", got.Total)
		}
	})

	t.Run("空の注文は400", func(t *testing.T) {
		rec := post(t, h, "/orders", `{"items":[]}`)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})
}

func TestCheckout(t *testing.T) {
	h := newTestHandler()

	t.Run("有効なカードで支払える", func(t *testing.T) {
		rec := post(t, h, "/checkout", `{
			"items": [{"sku":"MUG-01","quantity":1}],
			"card": {"number":"4242424242424242","exp_month":12,"exp_year":2027}
		}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body)
		}
		var got struct {
			Order     struct{ Total int } `json:"order"`
			ReceiptID string              `json:"receipt_id"`
		}
		decode(t, rec, &got)
		if got.Order.Total != 2198 {
			t.Errorf("order.Total = %d, want 2198", got.Order.Total)
		}
		if got.ReceiptID == "" {
			t.Error("receipt_id が空")
		}
	})
}

// newTestHandler は、2026年9月15日として動くハンドラーを返す。
func newTestHandler() http.Handler {
	cat := catalog.New(catalog.Product{SKU: "MUG-01", Name: "マグカップ", Price: 1999})
	gw := &payment.FakeGateway{Now: func() time.Time { return time.Date(2026, 9, 15, 10, 0, 0, 0, time.UTC) }}
	return NewHandler(cat, gw)
}

func post(t *testing.T, h http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("レスポンスを読めない: %v", err)
	}
}
