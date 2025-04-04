package main

import (
	"fmt"
	"github/vicpoo/APIproductora/handlers"
	"github/vicpoo/APIproductora/rabbitmq"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware para manejar CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Configurar RabbitMQ
	rabbitConn, err := rabbitmq.Connect("amqp://reyhades:reyhades@44.223.218.9:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Crear canal y configurar exchange/cola
	channel, err := rabbitmq.SetupRabbitMQ(rabbitConn)
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v", err)
	}
	defer channel.Close()

	// Configurar manejador HTTP con el canal RabbitMQ
	orderHandler := handlers.NewOrderHandler(channel)

	// Configurar rutas con Gorilla Mux
	router := mux.NewRouter()
	router.Use(corsMiddleware) // Aplicar middleware CORS
	router.HandleFunc("/api/order", orderHandler.HandleOrder).Methods("POST", "OPTIONS")

	// Iniciar servidor
	port := ":8004"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
