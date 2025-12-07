package logger

import (
	"encoding/json"
	"fmt"
	"go-heal/internal/types"
	"os"
	"sync"
)

func Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func Write(file *os.File, mu *sync.Mutex, entry types.LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Log error:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	file.Write(data)
	file.WriteString("\n")
}