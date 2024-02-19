package main

import (
    "encoding/json"
    "net/http"
)

type Notification struct {
    CustomerID string `json:"customer_id"`
    Message    string `json:"message"`
}

type NotificationRequest struct {
    CustomerID string `json:"customer_id"`
    Message    string `json:"message"`
}

func ReceiveNotificationEndpoint(w http.ResponseWriter, req *http.Request) {
    var notificationRequest NotificationRequest
    _ = json.NewDecoder(req.Body).Decode(&notificationRequest)

    err := notifyNotificationService(notificationRequest.CustomerID, notificationRequest.Message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(notificationRequest)
}

func notifyNotificationService(customerID, message string) error {
    // This is a placeholder for the actual implementation.
    // You would typically use a library or service like Firebase Cloud Messaging (FCM) to send the notification.
    // Here's an example of how you might do it with FCM:

    // data := map[string]string{
    // 	"customer_id": customerID,
    // 	"message":     message,
    // }

    // _, err := fcm.Send(fcm.Message{
    // 	To:   "/topics/all", // send to all devices subscribed to the "all" topic
    // 	Data: data,
    // })

    // return err

    return nil
}