package main

import (
	"log"
	"worker/internal/conversor"
)

func main() {
	err := conversor.ConvertFile("example-wpp-audio.ogg", "aac")
	if err != nil {
		log.Fatalf("Erro ao converter: %v", err)
	}
}
