package main

import (
    "database/sql"
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

// Remove the unused variable declaration

func GetOrderEndpoint(db *sql.DB, w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var order Order

    row := db.QueryRow("SELECT id, shopping_cart_id, customer_profile, product_location FROM orders WHERE id = ?", params["id"])
    err := row.Scan(&order.ID, &order.ShoppingCartID, &order.CustomerProfile, &order.ProductLocation)
    if err != nil {
        if err == sql.ErrNoRows {
            // No matching order found
            json.NewEncoder(w).Encode(&Order{})
            return
        }
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(order)
}


func GetAllOrdersEndpoint(db *sql.DB, w http.ResponseWriter, req *http.Request) {
    rows, err := db.Query("select id, shopping_cart_id, customer_profile, product_location from orders")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var orders []Order
    for rows.Next() {
        var order Order
        err = rows.Scan(&order.ID, &order.ShoppingCartID, &order.CustomerProfile, &order.ProductLocation)
        if err != nil {
            log.Fatal(err)
        }
        orders = append(orders, order)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(orders)
}
