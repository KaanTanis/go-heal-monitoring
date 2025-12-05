package main

import (
	"fmt"
	"net/http"
	"os"
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

func checkUrl(target TargetConfig) {
	for {
		processRequest(target)

		time.Sleep(target.Interval)
	}
}

func processRequest(target TargetConfig) {
	start := time.Now()
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(target.URL)
	if err != nil {
		elapsed := time.Since(start)

		fmt.Printf("[ERR] Target: %s - URL: %s - Error %s - Time: %v\n", target.Name, target.URL, err.Error(), elapsed)
		return
	}

	defer resp.Body.Close()
	
	elapsed := time.Since(start)

	status := "OK"
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = "WARN"
	}

	fmt.Printf("[LOG] Target: %s - URL: %s - Status: %s (%d) - Time: %v\n", target.Name, target.URL, status, resp.StatusCode, elapsed)
}

func main() {
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Fatal error loading config %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Agent started, loaded %d targets.\n", len(cfg.Targets))

	for _, target := range cfg.Targets {
		go checkUrl(target)
	}

	select {}
}