package models

type CodeRunner struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	UserID   uint   `json:"user_id"`
}

type CodeRunnerRes struct {
	CodeResult    string `json:"codeResult"`
	CorrelationID string `json:"correlationID"`
}
