package main

import (
	"context"
	"fmt"
	"go-heal/internal/config"
	"go-heal/internal/logger"
	"go-heal/internal/monitor"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}

	logFile, err := logger.Open("logs/monitor.log")
	if err != nil {
		fmt.Println("Log file can not be opened:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	var logMutex sync.Mutex

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	fmt.Printf("Agent started with %d targets. Press Ctrl+C to stop.\n", len(cfg.Targets))

	for _, target := range cfg.Targets {
		wg.Add(1)
		go monitor.CheckURL(ctx, &wg, target, logFile, &logMutex)
	}

	<-sigChan
	fmt.Println("\nShutdown signal received. Cleaning up...")

	cancel()
	wg.Wait()

	fmt.Println("All workers stopped. Bye!")
}