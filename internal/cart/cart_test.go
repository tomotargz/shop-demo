package cart

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("同じ商品は数量を足す", func(t *testing.T) {
		c := New()
		mustAdd(t, c, "MUG-01", 1)
		mustAdd(t, c, "MUG-01", 2)

		want := []Item{{SKU: "MUG-01", Quantity: 3}}
		if got := c.Items(); !reflect.DeepEqual(got, want) {
			t.Errorf("Items() = %v, want %v", got, want)
		}
	})

	t.Run("数量が0ならエラー", func(t *testing.T) {
		c := New()
		if err := c.Add("MUG-01", 0); err == nil {
			t.Error("Add() でエラーにならなかった")
		}
		if !c.IsEmpty() {
			t.Error("エラーのあとにカートが空でない")
		}
	})
}

func TestItems(t *testing.T) {
	c := New()
	mustAdd(t, c, "MUG-01", 1)
	mustAdd(t, c, "COFFEE-200", 2)

	want := []Item{{SKU: "COFFEE-200", Quantity: 2}, {SKU: "MUG-01", Quantity: 1}}
	if got := c.Items(); !reflect.DeepEqual(got, want) {
		t.Errorf("Items() = %v, want %v", got, want)
	}
}

func mustAdd(t *testing.T, c *Cart, sku string, quantity int) {
	t.Helper()
	if err := c.Add(sku, quantity); err != nil {
		t.Fatalf("Add(%q, %d): %v", sku, quantity, err)
	}
}
