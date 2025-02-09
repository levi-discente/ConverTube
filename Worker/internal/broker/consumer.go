package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"worker/internal/conversor"
	"worker/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(conn *amqp.Connection, reqQueue, resQueue string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		reqQueue,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	timeout := time.After(10 * time.Second)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				log.Println("Canal fechado, encerrando worker.")
				wg.Wait()
				return nil
			}
			wg.Add(1)
			go func(d amqp.Delivery) {
				defer wg.Done()

				var job ConversionJob
				if err := json.Unmarshal(d.Body, &job); err != nil {
					logError := "Invalid job format: " + err.Error()
					logger.SendLog(conn, "convert_log", job.OperationID, "error", logError)
					d.Nack(false, false)
					return
				}

				log.Printf("Processing job: %s", job.OperationID)
				logger.SendLog(conn, "convert_log", job.OperationID, "info", "Processing started")

				progressCallback := func(progress int) {
					response := ResponseMessage{
						OperationID: job.OperationID,
						Status:      "progress",
						Progress:    progress,
						Message:     "Processing...",
					}
					PublishResponse(conn, resQueue, response, d.CorrelationId)
					logger.SendLog(conn, "convert_log", job.OperationID, "info", fmt.Sprintf("Progress: %d%%", progress))
				}

				err := conversor.ConvertFile(job.FilePath, job.OutputFormat, job.Quality, progressCallback)
				if err != nil {
					logError := "Conversion failed: " + err.Error()
					logger.SendLog(conn, "convert_log", job.OperationID, "error", logError)

					response := ResponseMessage{
						OperationID: job.OperationID,
						Status:      "error",
						Message:     logError,
					}
					PublishResponse(conn, resQueue, response, d.CorrelationId)
					d.Nack(false, false)
					return
				}

				newFilePath := strings.TrimSuffix(job.FilePath, filepath.Ext(job.FilePath)) + "-converted." + job.OutputFormat
				logger.SendLog(conn, "convert_log", job.OperationID, "success", "Conversion completed successfully")

				response := ResponseMessage{
					OperationID: job.OperationID,
					Status:      "success",
					NewFilePath: newFilePath,
					Message:     "Conversion completed successfully!",
				}
				PublishResponse(conn, resQueue, response, d.CorrelationId)
				d.Ack(false)
			}(d)

		case <-timeout:
			log.Println("Timeout sem mensagens, encerrando worker.")
			wg.Wait()
			return nil
		}
	}
}
