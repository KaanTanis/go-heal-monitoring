package monitor

import (
	"context"
	"fmt"
	"go-heal/internal/logger"
	"go-heal/internal/types"
	"net/http"
	"os"
	"sync"
	"time"
)

func processRequest(target types.TargetConfig, file *os.File, mu *sync.Mutex) {
	start := time.Now()
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(target.URL)

	entry := types.LogEntry{
		TimeStamp:  time.Now().Format(time.RFC3339),
		TargetName: target.Name,
		URL:        target.URL,
	}

	if err != nil {
		entry.Level = "ERROR"
		entry.Error = err.Error()
		entry.DurationMs = time.Since(start).Milliseconds()

		logger.Write(file, mu, entry)

		fmt.Printf("[ERR] Target: %s - Error: %s\n", target.Name, err)
		return
	}

	defer resp.Body.Close()

	entry.Level = "INFO"
	entry.StatusCode = resp.StatusCode
	entry.DurationMs = time.Since(start).Milliseconds()

	status := "OK"
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		entry.Level = "WARN"
		status = "WARN"
	}

	logger.Write(file, mu, entry)

	fmt.Printf("[LOG] Target: %s - Status: %s (%d) - Time: %v\n", target.Name, status, resp.StatusCode, time.Since(start))
}

func CheckURL(ctx context.Context, wg *sync.WaitGroup, target types.TargetConfig, file *os.File, mu *sync.Mutex) {
	defer wg.Done()

	fmt.Printf("Started monitoring: %s\n", target.Name)

	ticker := time.NewTicker(target.Interval)
	defer ticker.Stop()

	// first control
	processRequest(target, file, mu)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Stopping monitor for: %s\n", target.Name)
			return
		case <-ticker.C:
			processRequest(target, file, mu)
		}
	}
}