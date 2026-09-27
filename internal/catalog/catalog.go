// Package catalog は、販売する商品の一覧を扱う。
package catalog

import (
	"fmt"

	"example.com/shop/internal/money"
)

// Product は販売する商品。Price は税抜きの単価。
type Product struct {
	SKU   string
	Name  string
	Price money.Yen
}

// Catalog は SKU から商品を引く。
type Catalog struct {
	products map[string]Product
}

// New は、渡された商品を持つカタログを返す。
func New(products ...Product) *Catalog {
	c := &Catalog{products: make(map[string]Product, len(products))}
	for _, p := range products {
		c.products[p.SKU] = p
	}
	return c
}

// Find は SKU に対応する商品を返す。見つからなければエラーを返す。
func (c *Catalog) Find(sku string) (Product, error) {
	p, ok := c.products[sku]
	if !ok {
		return Product{}, fmt.Errorf("商品が見つからない: %s", sku)
	}
	return p, nil
}

// Sample は、開発用の商品の一覧を返す。
func Sample() *Catalog {
	return New(
		Product{SKU: "COFFEE-200", Name: "ドリップコーヒー 200g", Price: 1280},
		Product{SKU: "MUG-01", Name: "マグカップ", Price: 1999},
		Product{SKU: "FILTER-100", Name: "ペーパーフィルター 100枚", Price: 398},
	)
}
