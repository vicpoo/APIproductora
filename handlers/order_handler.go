package handlers

import (
	"encoding/json"
	"github/vicpoo/APIproductora/models"
	"github/vicpoo/APIproductora/rabbitmq"
	"net/http"

	"github.com/streadway/amqp"
)

type OrderHandler struct {
	RabbitChannel *amqp.Channel
}

func NewOrderHandler(channel *amqp.Channel) *OrderHandler {
	return &OrderHandler{RabbitChannel: channel}
}

func (h *OrderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo de la solicitud
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if order.Name == "" || order.LastName == "" || order.Phone == "" || order.Email == "" || order.Quantity <= 0 || order.BedName == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Convertir a JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to process order", http.StatusInternalServerError)
		return
	}

	// Publicar en RabbitMQ
	err = rabbitmq.PublishOrder(h.RabbitChannel, orderJSON)
	if err != nil {
		http.Error(w, "Failed to send order to queue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order received successfully"})
}
