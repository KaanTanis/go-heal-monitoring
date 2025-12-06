package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

type TargetConfig struct {
	URL string `yaml:"url"`
	Name string `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
}

type Config struct {
	Targets []TargetConfig `yaml:"targets"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Config file could not be read: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("yaml parse error: %w", err)
	}

	return &cfg, nil
}

func processRequest(target TargetConfig) {
	start := time.Now()
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(target.URL)
	if err != nil {
		fmt.Printf("[ERR] Target: %s - Error: %s\n", target.Name, err)
		return
	}

	defer resp.Body.Close()

	status := "OK"
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = "WARN"
	}

	fmt.Printf("[LOG] Target: %s - Status: %s (%d) - Time: %v\n", target.Name, status, resp.StatusCode, time.Since(start))
}

func checkUrl(ctx context.Context, wg *sync.WaitGroup, target TargetConfig) {
	defer wg.Done()

	fmt.Printf("Started monitoring: %s\n", target.Name)

	ticker := time.NewTicker(target.Interval)
	defer ticker.Stop()

	// first control
	processRequest(target)

	for {
		select {
			case <-ctx.Done():
				fmt.Printf("Stopping monitor for: %s\n", target.Name)
				return
			case <-ticker.C:
				processRequest(target)
		}
	}
}

func main() {
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}

	// listen close channels
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// cancalable context
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	fmt.Printf("Agent started with %d targets. Press Ctrl+C to stop.\n", len(cfg.Targets))

	for _, target := range cfg.Targets {
		wg.Add(1)
		go checkUrl(ctx, &wg, target)
	}

	<-sigChan
	fmt.Println("\nShutdown signal received. Cleaning up...")

	cancel()

	fmt.Println("All workers stopped. Bye!")

}