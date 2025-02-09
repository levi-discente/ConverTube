package broker

import "time"

type ConversionJob struct {
	OperationID  string    `json:"operation_id"`
	FilePath     string    `json:"file_path"`
	FileName     string    `json:"file_name"`
	OutputFormat string    `json:"output_format"`
	RequestTime  time.Time `json:"request_time"`
	Quality      string    `json:"quality"`
}

type ConversionResult struct {
	OperationID string `json:"operation_id"`
	NewFilePath string `json:"new_file_path"`
	NewFileName string `json:"new_file_name"`
}

type ResponseMessage struct {
	OperationID string `json:"operation_id"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	Progress    int    `json:"progress,omitempty"`
	NewFilePath string `json:"new_file_path,omitempty"`
	NewFileName string `json:"new_file_name,omitempty"`
}
