package pkg

import (
	"fmt"
	"sync"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ConvertFile(inputPath, outputPath string) error {
	var wg sync.WaitGroup
	wg.Add(2)

	done := make(chan bool)

	go func() {
		defer wg.Done()
		showProgress(done)
	}()

	go func() {
		defer wg.Done()
		err := ffmpeg.Input(inputPath).
			Output(outputPath, ffmpeg.KwArgs{"c:v": "libx264", "preset": "fast", "crf": "23"}).
			OverWriteOutput().
			Run()
		if err != nil {
			fmt.Println("Erro na conversão:", err)
		}
		done <- true
	}()

	wg.Wait()
	close(done)
	return nil
}

func showProgress(done chan bool) {
	progress := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Println("\nConversão concluída!")
			return
		default:
			fmt.Printf("\rConvertendo... %s", progress[i%len(progress)])
			i++
			time.Sleep(500 * time.Millisecond)
		}
	}
}
