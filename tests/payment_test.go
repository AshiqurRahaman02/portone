package payment

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"portone/Controllers"
)

func TestCreateIntentHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/create_intent", payment.CreateIntentHandler).Methods("POST")

	envFilePath := filepath.Join("../../portone", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err)
	}

	requestBody := map[string]interface{}{
		"amount":   1000.00,
		"currency": "usd",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err, "Error creating test request body")

	req, err := http.NewRequest("POST", "/api/v1/create_intent", bytes.NewBuffer(requestBodyBytes))
	assert.NoError(t, err, "Error creating test request")

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error parsing JSON response")

	assert.True(t, response["success"].(bool), "Expected success to be true")
	assert.NotNil(t, response["payment_intent_id"], "Expected payment_intent_id to be present")
	assert.NotNil(t, response["payment"], "Expected payment to be present")

	paymentIntentID, ok := response["payment_intent_id"].(string)
	assert.True(t, ok, "Expected payment_intent_id to be a string")

	if paymentIntentID != "" {
		amount := response["payment"].(map[string]interface{})["amount"].(float64)
		assert.Equal(t, float64(1000), amount, "Unexpected payment amount")
	}
}

func TestConfirmIntentHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/confirm_payment_intent/{id}", payment.ConfirmIntentHandler).Methods("POST")

	ts := httptest.NewServer(r)
	defer ts.Close()

	envFilePath := filepath.Join("../../portone", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err)
	}

	// Create a new payment intent
	createIntentURL := ts.URL + "/api/v1/create_intent"
	createIntentReq, err := http.NewRequest("POST", createIntentURL, strings.NewReader(`{"amount": 2000, "currency": "usd"}`))
	if err != nil {
		t.Fatal(err)
	}

	createIntentRR := httptest.NewRecorder()
	createIntentHandler := http.HandlerFunc(payment.CreateIntentHandler)
	createIntentHandler.ServeHTTP(createIntentRR, createIntentReq)

	assert.Equal(t, http.StatusOK, createIntentRR.Code)

	var createIntentResponse map[string]interface{}
	err = json.Unmarshal(createIntentRR.Body.Bytes(), &createIntentResponse)
	if err != nil {
		t.Fatal(err)
	}

	// Use the created intent's ID for confirmation
	intentID := createIntentResponse["payment_intent_id"].(string)

	confirmIntentURL := ts.URL + "/api/v1/confirm_payment_intent/" + intentID
	confirmIntentReq, err := http.NewRequest("POST", confirmIntentURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": intentID}
	confirmIntentReq = mux.SetURLVars(confirmIntentReq, vars)

	confirmIntentRR := httptest.NewRecorder()
	confirmIntentHandler := http.HandlerFunc(payment.ConfirmIntentHandler)
	confirmIntentHandler.ServeHTTP(confirmIntentRR, confirmIntentReq)

	assert.Equal(t, http.StatusOK, confirmIntentRR.Code)

	var confirmIntentResponse map[string]interface{}
	err = json.Unmarshal(confirmIntentRR.Body.Bytes(), &confirmIntentResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, confirmIntentResponse["success"].(bool))
	assert.Equal(t, intentID, confirmIntentResponse["payment_intent_id"])
	assert.Equal(t, "requires_action", confirmIntentResponse["status"])
	assert.NotNil(t, confirmIntentResponse["client_secret"])
}

func TestGetIntentsHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/get_intents", payment.GetIntentsHandler).Methods("GET")

	envFilePath := filepath.Join("../../portone", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err)
	}

	req, err := http.NewRequest("GET", "/api/v1/get_intents", nil)
	assert.NoError(t, err, "Error creating test request")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error parsing JSON response")

	assert.True(t, response["success"].(bool), "Expected success to be true")
	assert.NotNil(t, response["paymentIntents"], "Expected paymentIntents to be present")

	paymentIntents, ok := response["paymentIntents"].([]interface{})
	assert.True(t, ok, "Expected paymentIntents to be a slice")

	if len(paymentIntents) > 0 {
		firstPaymentIntent := paymentIntents[0].(map[string]interface{})
		assert.NotNil(t, firstPaymentIntent["id"], "Expected paymentIntent ID to be present")
	}
}

func TestMain(m *testing.M) {
	m.Run()
}
