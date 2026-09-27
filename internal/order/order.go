// Package order は、カートの中身から注文を作る。
package order

import (
	"errors"

	"example.com/shop/internal/cart"
	"example.com/shop/internal/catalog"
	"example.com/shop/internal/money"
)

// ErrEmptyCart は、空のカートで注文しようとしたときのエラー。
var ErrEmptyCart = errors.New("カートが空")

// Line は注文の1行。
type Line struct {
	SKU       string
	Name      string
	UnitPrice money.Yen
	Quantity  int
	Amount    money.Yen
}

// Order は確定した注文の金額。
type Order struct {
	Lines    []Line
	Subtotal money.Yen // 税抜きの合計
	Tax      money.Yen
	Total    money.Yen // 税込みの支払額
}

// Place は、カートの中身から注文を作る。
func Place(c *cart.Cart, cat *catalog.Catalog) (Order, error) {
	if c.IsEmpty() {
		return Order{}, ErrEmptyCart
	}

	var o Order
	for _, item := range c.Items() {
		p, err := cat.Find(item.SKU)
		if err != nil {
			return Order{}, err
		}
		amount := p.Price * money.Yen(item.Quantity)
		o.Lines = append(o.Lines, Line{
			SKU:       p.SKU,
			Name:      p.Name,
			UnitPrice: p.Price,
			Quantity:  item.Quantity,
			Amount:    amount,
		})
		o.Subtotal += amount
	}
	o.Tax = money.Tax(o.Subtotal)
	o.Total = o.Subtotal + o.Tax
	return o, nil
}
