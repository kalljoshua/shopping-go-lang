package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NaySoftware/go-fcm"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *sql.DB

type Order struct {
	ID              string `json:"id,omitempty"`
	ShoppingCartID  string `json:"shopping_cart_id,omitempty"`
	CustomerProfile string `json:"customer_profile,omitempty"`
	ProductLocation string `json:"product_location,omitempty"`
}

type Notification struct {
	CustomerID string `json:"customer_id"`
	Message    string `json:"message"`
}

type ShippingRequest struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
}

type NotificationRequest struct {
	CustomerID string `json:"customer_id"`
	Message    string `json:"message"`
}

var orders []Order

func notifyNotificationService(customerID, message string) error {
	notificationRequest := &NotificationRequest{
		CustomerID: customerID,
		Message:    message,
	}

	requestBody, err := json.Marshal(notificationRequest)
	if err != nil {
		return err
	}

	_, err = http.Post("http://127.0.0.1:8080/notifications/send", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	return nil
}

func GetOrderEndpoint(w http.ResponseWriter, req *http.Request) {
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

func CreateOrderEndpoint(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	var order Order
	_ = json.NewDecoder(req.Body).Decode(&order)
	//order.ID = params["id"]
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

	// Notify the Notification Service about the new order
	err = notifyNotificationService(order.CustomerProfile, "Order created")
	if err != nil {
		log.Printf("Error notifying Notification Service: %v", err)
	}

	json.NewEncoder(w).Encode(order)
}

func GetAllOrdersEndpoint(w http.ResponseWriter, req *http.Request) {
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

func DeleteOrderEndpoint(w http.ResponseWriter, req *http.Request) {
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

	// Notify the Notification Service about the deleted order
	if (deletedOrder != Order{}) {
		err := notifyNotificationService(deletedOrder.CustomerProfile, "Order deleted")
		if err != nil {
			log.Printf("Error notifying Notification Service: %v", err)
		}
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

func ReceiveNotificationEndpoint(w http.ResponseWriter, req *http.Request) {
	var notification Notification
	_ = json.NewDecoder(req.Body).Decode(&notification)

	// Send a push notification
	var err error
	data := make(map[string]string)
	data["message"] = notification.Message
	data["customer_id"] = notification.CustomerID

	ids := []string{
		"device token or registration id", // replace with the device token or registration id
	}

	client := fcm.NewFcmClient("YourServerKey")
	client.NewFcmRegIdsMsg(ids, data)

	status, err := client.Send()

	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(notification)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
    create table if not exists orders (id text not null primary key, shopping_cart_id text, customer_profile text, product_location text);
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", GetOrderEndpoint).Methods("GET")
	router.HandleFunc("/orders/{id}", CreateOrderEndpoint).Methods("POST")
	router.HandleFunc("/orders/{id}", DeleteOrderEndpoint).Methods("DELETE")
	router.HandleFunc("/orders", GetAllOrdersEndpoint).Methods("GET")
	router.HandleFunc("/notifications/receive", ReceiveNotificationEndpoint).Methods("POST")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	router.HandleFunc("/ordersPage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ordersPage.html")
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}
