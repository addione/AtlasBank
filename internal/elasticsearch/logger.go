package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

// Logger represents an Elasticsearch logger
type Logger struct {
	client *elasticsearch.Client
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time              `json:"@timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// NewLogger creates a new Elasticsearch logger
func NewLogger(url string) *Logger {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error creating Elasticsearch client: %s\n", err)
		return &Logger{client: nil}
	}

	return &Logger{
		client: client,
	}
}

// Info logs an info message
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log("info", message, fields)
}

// Error logs an error message
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log("error", message, fields)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log("warn", message, fields)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log("debug", message, fields)
}

// log sends a log entry to Elasticsearch
func (l *Logger) log(level, message string, fields map[string]interface{}) {
	if l.client == nil {
		// Fallback to console logging if Elasticsearch is not available
		fmt.Printf("[%s] %s: %s %v\n", time.Now().Format(time.RFC3339), level, message, fields)
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("Error marshaling log entry: %s\n", err)
		return
	}

	// Index the log entry
	indexName := fmt.Sprintf("atlasbank-logs-%s", time.Now().Format("2006.01.02"))
	res, err := l.client.Index(
		indexName,
		bytes.NewReader(data),
		l.client.Index.WithContext(context.Background()),
	)
	if err != nil {
		fmt.Printf("Error indexing log entry: %s\n", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("Error response from Elasticsearch: %s\n", res.String())
	}
}

// Close closes the Elasticsearch logger
func (l *Logger) Close() {
	// Elasticsearch client doesn't need explicit closing
}
