package payment

import (
	"errors"
	"testing"
	"time"
)

func TestFakeGatewayCharge(t *testing.T) {
	// 2026年9月15日として動かす
	now := func() time.Time { return time.Date(2026, 9, 15, 10, 0, 0, 0, time.UTC) }

	tests := []struct {
		name     string
		card     Card
		wantCode string // 空なら成功
	}{
		{"有効期限内のカードは請求できる", Card{Number: "4242424242424242", ExpMonth: 12, ExpYear: 2027}, ""},
		{"期限の月の末日までは使える", Card{Number: "4242424242424242", ExpMonth: 9, ExpYear: 2026}, ""},
		{"期限の月を過ぎたカードは断る", Card{Number: "4242424242424242", ExpMonth: 8, ExpYear: 2026}, "expired_card"},
		{"期限の年を過ぎたカードは断る", Card{Number: "4242424242424242", ExpMonth: 12, ExpYear: 2025}, "expired_card"},
		{"0002で終わるカードは断る", Card{Number: "4000000000000002", ExpMonth: 12, ExpYear: 2027}, "card_declined"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &FakeGateway{Now: now}
			r, err := g.Charge(tt.card, 1000)

			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("Charge(): %v", err)
				}
				if r.Amount != 1000 {
					t.Errorf("Amount = %d, want 1000", r.Amount)
				}
				return
			}
			var de *DeclineError
			if !errors.As(err, &de) {
				t.Fatalf("err = %v, want DeclineError", err)
			}
			if de.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", de.Code, tt.wantCode)
			}
		})
	}
}
