package models

type CodeRunner struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type CodeRunnerRes struct {
	CodeResult    string `json:"codeResult"`
	CorrelationID string `json:"correlationID"`
}
