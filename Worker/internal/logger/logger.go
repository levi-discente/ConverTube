package logger

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendLog(conn *amqp.Connection, exchange string, operationID, level, message string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	logMsg := LogMessage{
		OperationID: operationID,
		Level:       level,
		Message:     message,
	}
	body, err := json.Marshal(logMsg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("[LOG] [%s] %s - %s", level, operationID, message)
	return nil
}
