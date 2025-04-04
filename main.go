package main

import (
	"fmt"
	"github/vicpoo/APIproductora/handlers"
	"github/vicpoo/APIproductora/rabbitmq"
	"log"
	"net/http"
)

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

	// Configurar rutas
	http.HandleFunc("/api/order", orderHandler.HandleOrder)

	// Iniciar servidor
	port := ":8004"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
