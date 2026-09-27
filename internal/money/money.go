// Package money は、金額（日本円）の計算を扱う。
package money

import "strconv"

// Yen は日本円の金額。1円未満は扱わない。
type Yen int64

// TaxRatePercent は消費税率（%）。
const TaxRatePercent = 10

// Tax は、税抜きの金額にかかる消費税額を返す。1円未満は切り捨てる。
func Tax(amount Yen) Yen {
	return amount * TaxRatePercent / 100
}

// String は「¥1,980」の形で金額を返す。
func (y Yen) String() string {
	s := strconv.FormatInt(int64(y), 10)
	sign := ""
	if y < 0 {
		sign, s = "-", s[1:]
	}
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return sign + "¥" + s
}
