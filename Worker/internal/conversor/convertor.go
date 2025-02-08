package conversor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var supportedFormats = map[string]bool{
	"mp4": true, "avi": true, "mkv": true, "mov": true,
	"flv": true, "webm": true, "ogg": true, "wav": true,
	"mp3": true, "aac": true, "flac": true, "wma": true,
	"gif": true,
}

func ConvertFile(inputPath, outputFormat string) error {
	if _, ok := supportedFormats[outputFormat]; !ok {
		return fmt.Errorf("formato de saída '%s' não suportado", outputFormat)
	}

	_, err := getFileFormat(inputPath)
	if err != nil {
		return fmt.Errorf("erro ao identificar o formato de entrada: %v", err)
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "-converted" + "." + outputFormat

	done := make(chan bool)
	progress := make(chan int, 1)

	go func() {
		if outputFormat == "gif" {
			err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": "1"}).
				Output(outputPath, ffmpeg.KwArgs{"s": "320x240", "pix_fmt": "rgb24", "t": "3", "r": "3"}).
				OverWriteOutput().
				WithOutput(nil, nil).
				Run()
			if err != nil {
				fmt.Println("Erro na conversão para GIF:", err)
			}
		} else {
			err := ffmpeg.Input(inputPath).
				Output(outputPath, ffmpeg.KwArgs{"c:v": "libx264", "preset": "fast", "crf": "23"}).
				OverWriteOutput().
				WithOutput(nil, nil).
				Run()
			if err != nil {
				fmt.Println("Erro na conversão:", err)
			}
		}
		done <- true
		close(done)
	}()

	go func() {
		showProgress(inputPath, outputPath, done, progress)
	}()

	for p := range progress {
		fmt.Printf("\rConvertendo... %d%%", p)
	}

	fmt.Println("\nConversão concluída!")
	return nil
}

func getFileFormat(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "format=format_name", "-of", "csv=p=0", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	format := strings.TrimSpace(string(output))
	return format, nil
}

func showProgress(inputPath, outputPath string, done chan bool, progress chan int) {
	defer close(progress)

	var lastSize int64 = 0
	for {
		select {
		case <-done:
			progress <- 100
			return
		default:
			time.Sleep(1 * time.Second)

			inputInfo, err := os.Stat(inputPath)
			if err != nil {
				continue
			}
			inputSize := inputInfo.Size()

			outputInfo, err := os.Stat(outputPath)
			if err != nil {
				continue
			}
			outputSize := outputInfo.Size()

			if inputSize > 0 {
				progressPercent := int(float64(outputSize) / float64(inputSize) * 100)
				if progressPercent > 100 {
					progressPercent = 100
				}
				if outputSize > lastSize {
					progress <- progressPercent
					lastSize = outputSize
				}
			}
		}
	}
}
