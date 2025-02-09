package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinIOStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no MinIO: %w", err)
	}

	return &MinIOStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *MinIOStorage) DownloadFile(objectName, localPath string) error {
	ctx := context.Background()

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo local: %w", err)
	}
	defer file.Close()

	obj, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("erro ao baixar arquivo do MinIO: %w", err)
	}
	defer obj.Close()

	_, err = io.Copy(file, obj)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo localmente: %w", err)
	}

	log.Printf("Arquivo %s baixado com sucesso para %s", objectName, localPath)
	return nil
}

func (s *MinIOStorage) UploadFile(localPath, objectName string) error {
	ctx := context.Background()

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo para upload: %w", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("erro ao obter informações do arquivo: %w", err)
	}

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("erro ao enviar arquivo para o MinIO: %w", err)
	}

	log.Printf("Arquivo %s enviado com sucesso para o MinIO", objectName)
	return nil
}
