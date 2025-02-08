package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Aviso: Arquivo .env não encontrado, usando variáveis de ambiente padrão")
	}

	return &Config{
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
	}
}
