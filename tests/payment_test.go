package payment

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

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

	intentID := "pi_3OcssqSJhfkaql8D04yP3W8X"
	requestURL := ts.URL + "/api/v1/confirm_payment_intent/" + intentID
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": intentID}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(payment.ConfirmIntentHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, response["success"].(bool))
	assert.Equal(t, intentID, response["payment_intent_id"])
	assert.Equal(t, "requires_action", response["status"])
	assert.NotNil(t, response["client_secret"])
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
