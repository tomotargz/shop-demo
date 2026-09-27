package order

import (
	"errors"
	"testing"

	"example.com/shop/internal/cart"
	"example.com/shop/internal/catalog"
	"example.com/shop/internal/money"
)

func TestPlace(t *testing.T) {
	cat := catalog.New(
		catalog.Product{SKU: "COFFEE-200", Name: "ドリップコーヒー 200g", Price: 1280},
		catalog.Product{SKU: "MUG-01", Name: "マグカップ", Price: 1999},
	)

	t.Run("1点の注文", func(t *testing.T) {
		c := cart.New()
		mustAdd(t, c, "COFFEE-200", 1)

		o := mustPlace(t, c, cat)
		assertAmounts(t, o, 1280, 128, 1408)
	})

	t.Run("複数の商品と数量", func(t *testing.T) {
		c := cart.New()
		mustAdd(t, c, "COFFEE-200", 2)
		mustAdd(t, c, "MUG-01", 1)

		o := mustPlace(t, c, cat)
		if len(o.Lines) != 2 {
			t.Fatalf("len(Lines) = %d, want 2", len(o.Lines))
		}
		if got := o.Lines[0].Amount; got != 2560 {
			t.Errorf("Lines[0].Amount = %d, want 2560", got)
		}
		assertAmounts(t, o, 4559, 455, 5014)
	})

	t.Run("消費税の1円未満は切り捨てる", func(t *testing.T) {
		c := cart.New()
		mustAdd(t, c, "MUG-01", 1)

		o := mustPlace(t, c, cat)
		assertAmounts(t, o, 1999, 199, 2198)
	})

	t.Run("空のカートは注文できない", func(t *testing.T) {
		_, err := Place(cart.New(), cat)
		if !errors.Is(err, ErrEmptyCart) {
			t.Errorf("err = %v, want ErrEmptyCart", err)
		}
	})

	t.Run("カタログにない商品は注文できない", func(t *testing.T) {
		c := cart.New()
		mustAdd(t, c, "UNKNOWN", 1)

		if _, err := Place(c, cat); err == nil {
			t.Error("Place() でエラーにならなかった")
		}
	})
}

func assertAmounts(t *testing.T, o Order, subtotal, tax, total money.Yen) {
	t.Helper()
	if o.Subtotal != subtotal {
		t.Errorf("Subtotal = %d, want %d", o.Subtotal, subtotal)
	}
	if o.Tax != tax {
		t.Errorf("Tax = %d, want %d", o.Tax, tax)
	}
	if o.Total != total {
		t.Errorf("Total = %d, want %d", o.Total, total)
	}
}

func mustAdd(t *testing.T, c *cart.Cart, sku string, quantity int) {
	t.Helper()
	if err := c.Add(sku, quantity); err != nil {
		t.Fatalf("Add(%q, %d): %v", sku, quantity, err)
	}
}

func mustPlace(t *testing.T, c *cart.Cart, cat *catalog.Catalog) Order {
	t.Helper()
	o, err := Place(c, cat)
	if err != nil {
		t.Fatalf("Place(): %v", err)
	}
	return o
}
