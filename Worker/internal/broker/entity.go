package broker

import "time"

type ConversionJob struct {
	OperationID  string    `json:"operation_id"`
	FilePath     string    `json:"file_path"`
	OutputFormat string    `json:"output_format"`
	RequestTime  time.Time `json:"request_time"`
	Quality      string    `json:"quality"`
}

type ConversionResult struct {
	OperationID string `json:"operation_id"`
	NewFilePath string `json:"new_file_path"`
}

type ResponseMessage struct {
	OperationID string `json:"operation_id"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	Progress    int    `json:"progress,omitempty"`
	NewFilePath string `json:"new_file_path,omitempty"`
}
