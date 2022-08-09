package worker

import "encoding/json"

type CodeLine struct {
	Content string `json:"content"`
	Line    int    `json:"line"`
}

type BacktraceFile struct {
	File       string     `json:"file"`
	Line       int        `json:"line"`
	Column     int        `json:"column"`
	Function   *string    `json:"function"`
	SourceCode []CodeLine `json:"sourceCode"`
}

type Payload struct {
	Title          string           `json:"title"`
	Type           string           `json:"type"`
	Backtrace      []BacktraceFile  `json:"backtrace"`
	Context        *json.RawMessage `json:"context"`
	CatcherVersion string           `json:"catcherVersion"`
	Timestamp      int              `json:"timestamp"`
}

type Event struct {
	ProjectId   string  `json:"projectId"`
	Payload     Payload `json:"payload"`
	CatcherType string  `json:"catcherType"`
}
