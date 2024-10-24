package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

type OrderDetails struct {
	OrderAmount     float64  `json:"order_amount"`
	OrderCurrency   string   `json:"order_currency"`
	OrderID         string   `json:"order_id"`
	CustomerDetails Customer `json:"customer_details"`
	OrderMeta       Meta     `json:"order_meta"`
}

type Customer struct {
	CustomerID    string `json:"customer_id"`
	CustomerPhone string `json:"customer_phone"`
}

type Meta struct {
	ReturnURL string `json:"return_url"`
}

func GenrandOrderId() string {
	return "devstudio_" + strconv.Itoa(rand.Intn(10000000))
}

func CashHandler(w http.ResponseWriter, r *http.Request) {
	order := OrderDetails{
		OrderAmount:   200.00,
		OrderCurrency: "INR",
		OrderID:       GenrandOrderId(),
		CustomerDetails: Customer{
			CustomerID:    "devstudio_user",
			CustomerPhone: "8474090589",
		},
		OrderMeta: Meta{
			ReturnURL: "https://www.cashfree.com/devstudio/preview/pg/web/checkout?order_id={order_id}",
		},
	}

	// Convert order details to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to marshal order details", http.StatusInternalServerError)
		return
	}

	// Prepare the request
	url := "https://sandbox.cashfree.com/pg/orders"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(orderJSON))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("x-client-id", "TEST430329ae80e0f32e41a393d78b923034")
	req.Header.Set("x-client-secret", "TESTaf195616268bd6202eeb3bf8dc458956e7192a85")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-version", "2023-08-01")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
}
