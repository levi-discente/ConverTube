package tests

import (
	"os"
	"testing"
	"worker/internal/conversor"
)

func TestConvertFile(t *testing.T) {
	inputFile := "../../example-full-hd.mp4"
	outputFormat := "mp4"
	quality := "low"

	expectedOutput := "../../example-full-hd-converted." + outputFormat

	err := conversor.ConvertFile(inputFile, outputFormat, quality, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if _, err := os.Stat(expectedOutput); os.IsNotExist(err) {
		t.Fatalf("Converted file not found: %s", expectedOutput)
	}

	_ = os.Remove(expectedOutput)
}
