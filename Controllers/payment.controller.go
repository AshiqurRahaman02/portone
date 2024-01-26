package payment

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
	"net/http"
	"os"
)

func CreateIntentHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to parse request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	amountFloat, ok := requestBody["amount"].(float64)
	if !ok {
		http.Error(w, "Missing or invalid 'amount' parameter", http.StatusBadRequest)
		return
	}
	currency, ok := requestBody["currency"].(string)
	if !ok {
		http.Error(w, "Missing or invalid 'currency' parameter", http.StatusBadRequest)
		return
	}
	amount := int64(amountFloat)

	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		fmt.Println("STRIPE_API_KEY is not set. Please set it.")
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}
	stripe.Key = apiKey

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amount),
		Currency:      stripe.String(currency),
		PaymentMethod: stripe.String("pm_card_visa"),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create PaymentIntent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":           true,
		"payment_intent_id": pi.ID,
		"payment":           pi,
	}
	json.NewEncoder(w).Encode(response)
}

func ConfirmIntentHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		fmt.Println("STRIPE_API_KEY is not set. Please set it.")
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}
	stripe.Key = apiKey

	vars := mux.Vars(r)
	intentID := vars["id"]

	intent, err := paymentintent.Get(intentID, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve PaymentIntent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if intent.Status != "requires_confirmation" {
		errorMessage := fmt.Sprintf("PaymentIntent cannot be captured in its current state, State: %s", intent.Status)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	confirmParams := &stripe.PaymentIntentConfirmParams{}
	confirmedIntent, err := paymentintent.Confirm(intent.ID, confirmParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to confirm PaymentIntent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":           true,
		"payment_intent_id": confirmedIntent.ID,
		"status":            confirmedIntent.Status,
		"client_secret":     confirmedIntent.ClientSecret,
	}
	json.NewEncoder(w).Encode(response)
}

func CaptureIntentHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		fmt.Println("STRIPE_API_KEY is not set. Please set it.")
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}
	stripe.Key = apiKey

	vars := mux.Vars(r)
	intentID := vars["id"]

	intent, err := paymentintent.Get(intentID, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve PaymentIntent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if intent.Status != "requires_capture" {
		errorMessage := fmt.Sprintf("PaymentIntent cannot be captured in its current state, State: %s", intent.Status)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	params := &stripe.PaymentIntentCaptureParams{}
	intent, err = paymentintent.Capture(intent.ID, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to capture PaymentIntent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":           true,
		"payment_intent_id": intent.ID,
		"status":            intent.Status,
	}
	json.NewEncoder(w).Encode(response)
}

func CreateRefundHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		fmt.Println("STRIPE_API_KEY is not set. Please set it.")
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}
	stripe.Key = apiKey

	vars := mux.Vars(r)
	intentID := vars["id"]

	intent, err := paymentintent.Get(intentID, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve PaymentIntent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if intent.Status != "succeeded" {
		errorMessage := fmt.Sprintf("PaymentIntent cannot be refunded in its current state, State: %s", intent.Status)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(intentID),
	}
	refund, err := refund.New(refundParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create refund: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":         true,
		"refund_id":       refund.ID,
		"status":          refund.Status,
		"amount_refunded": refund.Amount,
	}
	json.NewEncoder(w).Encode(response)
}

func GetIntentsHandler(w http.ResponseWriter, r *http.Request) {

	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		fmt.Println("STRIPE_API_KEY is not set. Please set it.")
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}
	stripe.Key = apiKey

	params := &stripe.PaymentIntentListParams{}
	intents := paymentintent.List(params)

	var intentList []map[string]interface{}
	for intents.Next() {
		intent := intents.PaymentIntent()
		intentList = append(intentList, map[string]interface{}{
			"id":           intent.ID,
			"amount":       intent.Amount,
			"currency":     intent.Currency,
			"status":       intent.Status,
			"created":      intent.Created,
			"clientSecret": intent.ClientSecret,
		})
	}

	response := map[string]interface{}{
		"success":        true,
		"paymentIntents": intentList,
	}
	json.NewEncoder(w).Encode(response)
}
