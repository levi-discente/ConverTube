package main

import (
	"log"
	"myproject/pkg/ffmpeg"
)

func main() {
	err := ffmpeg.ConvertFile("input.mp4", "output.avi")
	if err != nil {
		log.Fatalf("Erro ao converter: %v", err)
	}
}
