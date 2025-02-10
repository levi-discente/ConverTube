package broker

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"worker/internal/conversor"
	"worker/internal/logger"
	"worker/internal/storage"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(conn *amqp.Connection, reqQueue, resQueue string) error {
	storageClient, err := storage.NewMinIOStorage(
		"minio-headless.default.svc.cluster.local:9000",
		"minio",
		"minio123",
		"uploads",
		false,
	)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MinIO: %v", err)
	}

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
	timeout := time.After(30 * time.Second)

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
					sendError(conn, resQueue, d, job.OperationID, "Invalid job format: "+err.Error())
					return
				}

				log.Printf("Processing job: %s", job.OperationID)
				logger.SendLog(conn, "convert_log", job.OperationID, "info", "Processing started")
				log.Printf("Job recebido: %+v", job)

				localFilePath := "/tmp/" + job.FileName
				defer os.Remove(localFilePath)

				if err := storageClient.DownloadFile(job.FilePath, localFilePath); err != nil {
					sendError(conn, resQueue, d, job.OperationID, "Failed to download file from MinIO: "+err.Error())
					return
				}

				progressCallback := func(progress int) {
					sendResponse(conn, resQueue, d.CorrelationId, job.OperationID, "progress", progress, "Processing...")
				}

				if err := conversor.ConvertFile(localFilePath, job.OutputFormat, job.Quality, progressCallback); err != nil {
					sendError(conn, resQueue, d, job.OperationID, "Conversion failed: "+err.Error())
					return
				}

				newFileName := strings.TrimSuffix(job.FileName, filepath.Ext(job.FileName)) + "-converted." + job.OutputFormat
				newFilePath := "/tmp/" + newFileName
				defer os.Remove(newFilePath)

				if err := storageClient.UploadFile(newFilePath, newFileName); err != nil {
					sendError(conn, resQueue, d, job.OperationID, "Failed to upload converted file to MinIO: "+err.Error())
					return
				}

				logger.SendLog(conn, "convert_log", job.OperationID, "success", "Conversion completed successfully")
				sendResponse(conn, resQueue, d.CorrelationId, job.OperationID, "success", 100, "Conversion completed successfully!", newFileName)
				d.Ack(false)
			}(d)

		case <-timeout:
			log.Println("Timeout sem mensagens, encerrando worker.")
			wg.Wait()
			return nil
		}
	}
}

func sendResponse(conn *amqp.Connection, resQueue, correlationID, operationID, status string, progress int, message string, fileName ...string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Erro ao abrir canal para resposta: %v", err)
		return
	}
	defer ch.Close()

	response := ResponseMessage{
		OperationID: operationID,
		Status:      status,
		Progress:    progress,
		Message:     message,
	}
	if len(fileName) > 0 {
		response.NewFileName = fileName[0]
	}

	body, _ := json.Marshal(response)
	err = ch.Publish(
		"",
		resQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			Body:          body,
		},
	)
	if err != nil {
		log.Printf("Erro ao publicar resposta para RabbitMQ: %v", err)
	}
}

func sendError(conn *amqp.Connection, resQueue string, d amqp.Delivery, operationID, errorMessage string) {
	logger.SendLog(conn, "convert_log", operationID, "error", errorMessage)
	sendResponse(conn, resQueue, d.CorrelationId, operationID, "error", 0, errorMessage)
	d.Nack(false, false)
}
