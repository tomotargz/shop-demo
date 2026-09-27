// shop は、注文と決済を受け付けるHTTPサーバー。
package main

import (
	"log"
	"net/http"
	"time"

	"example.com/shop/internal/api"
	"example.com/shop/internal/catalog"
	"example.com/shop/internal/payment"
)

func main() {
	h := api.NewHandler(catalog.Sample(), &payment.FakeGateway{Now: time.Now})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
