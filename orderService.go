package main

import (
    "database/sql"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

type Order struct {
    ID              string `json:"id,omitempty"`
    ShoppingCartID  string `json:"shopping_cart_id,omitempty"`
    CustomerProfile string `json:"customer_profile,omitempty"`
    ProductLocation string `json:"product_location,omitempty"`
}

// Remove the unused variable declaration

func CreateOrderEndpoint(db *sql.DB, w http.ResponseWriter, req *http.Request) {
    var order Order
    _ = json.NewDecoder(req.Body).Decode(&order)
    order.ID = uuid.New().String()

    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    stmt, err := tx.Prepare("insert into orders(id, shopping_cart_id, customer_profile, product_location) values(?, ?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    _, err = stmt.Exec(order.ID, order.ShoppingCartID, order.CustomerProfile, order.ProductLocation)
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()

    json.NewEncoder(w).Encode(order)
}


func DeleteOrderEndpoint(db *sql.DB, w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var deletedOrder Order

    // Fetch the order to be deleted
    row := db.QueryRow("SELECT id, shopping_cart_id, customer_profile, product_location FROM orders WHERE id = ?", params["id"])
    err := row.Scan(&deletedOrder.ID, &deletedOrder.ShoppingCartID, &deletedOrder.CustomerProfile, &deletedOrder.ProductLocation)
    if err != nil {
        if err == sql.ErrNoRows {
            // No matching order found
            json.NewEncoder(w).Encode(&Order{})
            return
        }
        log.Fatal(err)
    }

    // Delete the order
    _, err = db.Exec("DELETE FROM orders WHERE id = ?", params["id"])
    if err != nil {
        log.Fatal(err)
    }

    // Fetch all remaining orders
    rows, err := db.Query("SELECT id, shopping_cart_id, customer_profile, product_location FROM orders")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var orders []Order
    for rows.Next() {
        var order Order
        if err := rows.Scan(&order.ID, &order.ShoppingCartID, &order.CustomerProfile, &order.ProductLocation); err != nil {
            log.Fatal(err)
        }
        orders = append(orders, order)
    }

    json.NewEncoder(w).Encode(orders)
}