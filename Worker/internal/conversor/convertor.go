package conversor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var (
	supportedFormats = map[string]bool{
		"mp4": true, "avi": true, "mkv": true, "mov": true,
		"flv": true, "webm": true, "ogg": true, "wav": true,
		"mp3": true, "aac": true, "flac": true, "wma": true,
		"gif": true,
	}

	qualityPresets = map[string]struct {
		preset string
		crf    int
	}{
		"low":             {"ultrafast", 30},
		"medium":          {"fast", 23},
		"high":            {"slow", 18},
		"max_compression": {"veryslow", 26},
	}
)

func ConvertFile(inputPath, outputFormat, quality string, progressCallback func(progress int)) error {
	if !isFormatSupported(outputFormat) {
		return fmt.Errorf("output format '%s' is not supported", outputFormat)
	}

	q, ok := qualityPresets[quality]
	if !ok {
		return fmt.Errorf("invalid quality setting '%s'", quality)
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "-converted." + outputFormat

	var wg sync.WaitGroup
	done := make(chan bool)
	progress := make(chan int, 1)
	errChan := make(chan error, 1)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := processConversion(inputPath, outputPath, outputFormat, q.preset, q.crf)
		if err != nil {
			errChan <- err
		}
		done <- true
	}()

	go func() {
		defer wg.Done()
		trackProgress(inputPath, outputPath, done, progress)
	}()

	for p := range progress {
		if progressCallback != nil {
			progressCallback(p)
		}
	}

	wg.Wait()
	close(errChan)

	if err, ok := <-errChan; ok {
		return err
	}

	fmt.Println("\nConversion completed successfully!")
	return nil
}

func isFormatSupported(format string) bool {
	_, exists := supportedFormats[format]
	return exists
}

func processConversion(inputPath, outputPath, format, preset string, crf int) error {
	var err error
	if format == "gif" {
		err = ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": "1"}).
			Output(outputPath, ffmpeg.KwArgs{"s": "320x240", "pix_fmt": "rgb24", "t": "3", "r": "3"}).
			OverWriteOutput().
			WithOutput(nil, nil).
			Run()
	} else {
		err = ffmpeg.Input(inputPath).
			Output(outputPath, ffmpeg.KwArgs{"c:v": "libx264", "preset": preset, "crf": fmt.Sprintf("%d", crf)}).
			OverWriteOutput().
			WithOutput(nil, nil).
			Run()
	}
	return err
}

func trackProgress(inputPath, outputPath string, done chan bool, progress chan int) {
	defer close(progress)

	var lastSize int64
	for {
		select {
		case <-done:
			progress <- 100
			return
		default:
			time.Sleep(1 * time.Second)

			outputSize, err := getFileSize(outputPath)
			if err != nil {
				continue
			}

			inputSize, err := getFileSize(inputPath)
			if err != nil || inputSize == 0 {
				continue
			}

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

func getFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
