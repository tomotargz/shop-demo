package payment

import (
	"fmt"
	"strings"
	"time"

	"example.com/shop/internal/money"
)

// FakeGateway は、決済代行の動きをまねる。本物の決済代行と契約するまでの間、開発とテストで使う。
//
//   - 有効期限切れのカードは expired_card で断る。カードは有効期限の月の末日まで使える
//   - 番号が 0002 で終わるカードは card_declined で断る
type FakeGateway struct {
	Now func() time.Time // 現在時刻。テストでは固定した時刻を返す関数を渡す

	charges int
}

// Charge はカードに amount を請求する。
func (g *FakeGateway) Charge(card Card, amount money.Yen) (Receipt, error) {
	now := g.Now()
	if card.ExpYear < now.Year() || (card.ExpYear == now.Year() && card.ExpMonth < int(now.Month())) {
		return Receipt{}, &DeclineError{Code: "expired_card"}
	}
	if strings.HasSuffix(card.Number, "0002") {
		return Receipt{}, &DeclineError{Code: "card_declined"}
	}
	g.charges++
	return Receipt{ID: fmt.Sprintf("ch_%04d", g.charges), Amount: amount}, nil
}
