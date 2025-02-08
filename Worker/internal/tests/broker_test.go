package tests

import (
	"encoding/json"
	"testing"
	"time"
	"worker/internal/broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMQCommunication(t *testing.T) {
	// 🔥 Conectar ao RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 🔥 Criar canal
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 🔥 Declarar fila temporária
	queue, err := ch.QueueDeclare(
		"test_queue",
		false,
		true, // 🔥 Auto-delete ativado para limpar após os testes
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	// 🔥 Criar uma mensagem de teste
	message := broker.ResponseMessage{
		OperationID: "test123",
		Status:      "success",
		Message:     "Test message",
	}

	// 🔥 Serializar a mensagem para JSON
	body, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// 🔥 Publicar a mensagem
	err = ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// 🔥 Consumir a mensagem
	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to consume message: %v", err)
	}

	// 🔥 Ler a mensagem dentro de 2 segundos
	select {
	case msg := <-msgs:
		var receivedMsg broker.ResponseMessage
		if err := json.Unmarshal(msg.Body, &receivedMsg); err != nil {
			t.Fatalf("Failed to unmarshal received message: %v", err)
		}

		// 🔥 Verificar se a mensagem recebida é a esperada
		if receivedMsg.OperationID != message.OperationID {
			t.Fatalf("Unexpected OperationID: got %s, expected %s", receivedMsg.OperationID, message.OperationID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: Did not receive message")
	}
}
