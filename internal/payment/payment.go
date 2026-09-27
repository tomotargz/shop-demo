// Package payment は、カードでの支払いを扱う。
package payment

import "example.com/shop/internal/money"

// Card は支払いに使うカード。
type Card struct {
	Number   string
	ExpMonth int // 有効期限の月（1〜12）
	ExpYear  int // 有効期限の年（西暦4桁）
}

// Receipt は、決済が成功したときの控え。
type Receipt struct {
	ID     string
	Amount money.Yen
}

// Gateway は、カードに請求する決済代行サービス。
type Gateway interface {
	Charge(card Card, amount money.Yen) (Receipt, error)
}

// DeclineError は、決済代行がカードを受け付けなかったときのエラー。
type DeclineError struct {
	Code string // 決済代行が返す理由。例：expired_card、card_declined
}

func (e *DeclineError) Error() string {
	return "カードが受け付けられなかった: " + e.Code
}
