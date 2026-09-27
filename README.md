# shop

コーヒー豆と器具を売る通販サイトの、注文と決済を受け付けるサーバー。
カートに入った商品から注文を作り、カードで支払ってもらう。

## 構成

- cmd/shop           HTTPサーバーを起動する
- internal/api       注文と決済のAPI
- internal/cart      カート（商品と数量）
- internal/catalog   販売する商品の一覧
- internal/money     金額と消費税の計算
- internal/order     カートから注文を作る
- internal/payment   カードでの支払い。決済代行はまだ偽物

## 使い方

テストを実行する

    go test ./...

サーバーを起動して、支払う

    go run ./cmd/shop
    curl -X POST localhost:8080/checkout -d '{
      "items": [{"sku": "MUG-01", "quantity": 1}],
      "card": {"number": "4242424242424242",
               "exp_month": 12, "exp_year": 2030}}'
