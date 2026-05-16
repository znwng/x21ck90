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

func Log(method string, statusCode int, path string, clientIP string) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Method:    method,
		Status:    statusCode,
		ClientIP:  clientIP,
		Path:      path,
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	file, err := os.OpenFile(
		filepath.Join(wd, "log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(entry)
	if err != nil {
		return
	}
}

