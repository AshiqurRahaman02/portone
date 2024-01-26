package routes

import (
	"github.com/gorilla/mux"
	"portone/Controllers"
)

func SetupRoutes() *mux.Router {
	paymentRouter := mux.NewRouter()

	paymentRouter.HandleFunc("/api/v1/create_intent", payment.CreateIntentHandler).Methods("POST")
	paymentRouter.HandleFunc("/api/v1/confirm_payment_intent/{id}", payment.ConfirmIntentHandler).Methods("POST")
	paymentRouter.HandleFunc("/api/v1/capture_intent/{id}", payment.CaptureIntentHandler).Methods("POST")
	paymentRouter.HandleFunc("/api/v1/create_refund/{id}", payment.CreateRefundHandler).Methods("POST")
	paymentRouter.HandleFunc("/api/v1/get_intents", payment.GetIntentsHandler).Methods("GET")

	return paymentRouter
}
