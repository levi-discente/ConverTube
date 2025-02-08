package main

import (
	"log"
	"os"
	"worker/config"
	"worker/internal/broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/main.go <queue_id>")
	}
	queueID := os.Args[1]
	reqQueue := "conversion_jobs_" + queueID
	resQueue := "conversion_responses_" + queueID

	cfg := config.LoadConfig()
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Printf("Worker iniciado para a fila: %s", reqQueue)
	if err := broker.StartConsumer(conn, reqQueue, resQueue); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
