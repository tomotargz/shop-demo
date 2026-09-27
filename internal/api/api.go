// Package api は、注文と決済を受け付けるHTTPのAPI。
package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"example.com/shop/internal/cart"
	"example.com/shop/internal/catalog"
	"example.com/shop/internal/order"
	"example.com/shop/internal/payment"
)

type server struct {
	catalog *catalog.Catalog
	gateway payment.Gateway
}

// NewHandler は、APIのハンドラーを返す。
func NewHandler(cat *catalog.Catalog, gw payment.Gateway) http.Handler {
	s := &server{catalog: cat, gateway: gw}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", s.handleOrders)
	mux.HandleFunc("POST /checkout", s.handleCheckout)
	return mux
}

type itemRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type orderRequest struct {
	Items []itemRequest `json:"items"`
}

type checkoutRequest struct {
	Items []itemRequest `json:"items"`
	Card  struct {
		Number   string `json:"number"`
		ExpMonth int    `json:"exp_month"`
		ExpYear  int    `json:"exp_year"`
	} `json:"card"`
}

type checkoutResponse struct {
	Order     order.Order `json:"order"`
	ReceiptID string      `json:"receipt_id"`
}

// handleOrders は、注文の金額を計算して返す。
func (s *server) handleOrders(w http.ResponseWriter, r *http.Request) {
	var req orderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストの形式が正しくない", http.StatusBadRequest)
		return
	}
	o, ok := s.placeOrder(w, req.Items)
	if !ok {
		return
	}
	writeJSON(w, o)
}

// handleCheckout は、注文を確定してカードで支払う。
func (s *server) handleCheckout(w http.ResponseWriter, r *http.Request) {
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストの形式が正しくない", http.StatusBadRequest)
		return
	}
	o, ok := s.placeOrder(w, req.Items)
	if !ok {
		return
	}

	card := payment.Card{Number: req.Card.Number, ExpMonth: req.Card.ExpMonth, ExpYear: req.Card.ExpYear}
	receipt, err := s.gateway.Charge(card, o.Total)
	if err != nil {
		log.Printf("決済に失敗した: %v", err)
		http.Error(w, "システムエラーが発生しました", http.StatusInternalServerError)
		return
	}
	writeJSON(w, checkoutResponse{Order: o, ReceiptID: receipt.ID})
}

// placeOrder は、リクエストの商品から注文を作る。作れなかったときはエラーを書き込み、false を返す。
func (s *server) placeOrder(w http.ResponseWriter, items []itemRequest) (order.Order, bool) {
	c := cart.New()
	for _, it := range items {
		if err := c.Add(it.SKU, it.Quantity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return order.Order{}, false
		}
	}

	o, err := order.Place(c, s.catalog)
	if errors.Is(err, order.ErrEmptyCart) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return order.Order{}, false
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return order.Order{}, false
	}
	return o, true
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("レスポンスを書けなかった: %v", err)
	}
}
