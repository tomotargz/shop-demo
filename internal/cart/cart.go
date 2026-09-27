// Package cart は、注文前に商品を入れておくカートを扱う。
package cart

import (
	"errors"
	"sort"
)

// Cart は、SKU ごとの数量を持つ。
type Cart struct {
	items map[string]int
}

// Item はカートの1行。
type Item struct {
	SKU      string
	Quantity int
}

// New は空のカートを返す。
func New() *Cart {
	return &Cart{items: make(map[string]int)}
}

// Add は商品を quantity 個追加する。すでにある商品なら数量を足す。
func (c *Cart) Add(sku string, quantity int) error {
	if quantity <= 0 {
		return errors.New("数量は1以上にする")
	}
	c.items[sku] += quantity
	return nil
}

// Items は、SKU の順に並べたカートの中身を返す。
func (c *Cart) Items() []Item {
	items := make([]Item, 0, len(c.items))
	for sku, q := range c.items {
		items = append(items, Item{SKU: sku, Quantity: q})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].SKU < items[j].SKU })
	return items
}

// IsEmpty は、カートに何も入っていなければ true を返す。
func (c *Cart) IsEmpty() bool {
	return len(c.items) == 0
}
