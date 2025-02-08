package tests

import (
	"encoding/json"
	"testing"
	"time"
	"worker/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestLogExchange(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"convert_log",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare test queue: %v", err)
	}
	err = ch.QueueBind(q.Name, "", "convert_log", false, nil)
	if err != nil {
		t.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	err = logger.SendLog(conn, "convert_log", "test_operation", "info", "Test log message")
	if err != nil {
		t.Fatalf("Failed to send log message: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to consume messages: %v", err)
	}

	select {
	case msg := <-msgs:
		var logMsg logger.LogMessage
		if err := json.Unmarshal(msg.Body, &logMsg); err != nil {
			t.Fatalf("Failed to unmarshal log message: %v", err)
		}

		if logMsg.Message != "Test log message" {
			t.Fatalf("Unexpected log message: got %s, expected %s", logMsg.Message, "Test log message")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: Did not receive log message")
	}
}
