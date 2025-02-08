package logger

type LogMessage struct {
	OperationID string `json:"operation_id"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}
