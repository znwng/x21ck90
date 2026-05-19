package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Status    int    `json:"status"`
	ClientIP  string `json:"client_ip"`
	Path      string `json:"path"`
}

func OpenOrCreate(name string) (*os.File, error) {
	return os.OpenFile(
		name,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
}

func Log(method string, statusCode int, path string, clientIP string) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Method:    method,
		Status:    statusCode,
		ClientIP:  clientIP,
		Path:      path,
	}

	logDirPath := os.Getenv("LOG_DIR")

	filename := filepath.Join(
		logDirPath,
		time.Now().Format("02012006"),
	)

	file, err := OpenOrCreate(filename)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(entry); err != nil {
		return
	}
}

