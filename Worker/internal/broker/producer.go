package broker

import (
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishResponse(conn *amqp.Connection, resQueue string, response ResponseMessage, corrID string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		resQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: corrID,
		},
	)
	return err
}

func PublishProgress(conn *amqp.Connection, progressQueue string, operationID string, progress int, corrID string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		progressQueue,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msg := map[string]interface{}{
		"operation_id": operationID,
		"progress":     progress,
		"timestamp":    time.Now().Format(time.RFC3339),
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: corrID,
		},
	)
	return err
}

func PublishError(conn *amqp.Connection, errorQueue string, errorMsg string, corrID string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		errorQueue,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msg := map[string]interface{}{
		"error":     errorMsg,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: corrID,
		},
	)
	return err
}

func PublishLog(conn *amqp.Connection, logQueue string, response ResponseMessage, corrID string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		logQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: corrID,
		},
	)
	return err
}
