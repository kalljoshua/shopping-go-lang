package main

/* import (
    "bytes"
    "encoding/json"
    "net/http"
)

type ShippingRequest struct {
    OrderID    string `json:"order_id"`
    CustomerID string `json:"customer_id"`
}

type NotificationRequest struct {
    CustomerID string `json:"customer_id"`
    Message    string `json:"message"`
} */

/* func alertShippingService(orderID, customerID string) error {
    shippingRequest := &ShippingRequest{
        OrderID:    orderID,
        CustomerID: customerID,
    }

    requestBody, err := json.Marshal(shippingRequest)
    if err != nil {
        return err
    }

    _, err = http.Post("https://127.0.0.1:8080/api/shipping/receive", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }

    return nil
} */

/* func notifyNotificationService(customerID, message string) error {
    notificationRequest := &NotificationRequest{
        CustomerID: customerID,
        Message:    message,
    }

    requestBody, err := json.Marshal(notificationRequest)
    if err != nil {
        return err
    }

    _, err = http.Post("https://127.0.0.1:8081/notifications/send", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }

    return nil
} */