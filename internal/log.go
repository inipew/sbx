package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchLog monitors the specified log file for changes and prints new log entries.
func WatchLog(ctx context.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		return fmt.Errorf("failed to add file to watcher: %w", err)
	}

	reader := bufio.NewReader(file)
	file.Seek(0, io.SeekEnd)

	logChan := make(chan string)
	done := make(chan struct{})

	// Print new log entries
	go func() {
		defer close(done)
		for line := range logChan {
			fmt.Print(line)
		}
	}()

	// Monitor file changes
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Context cancelled, exit goroutine
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					for {
						line, err := reader.ReadString('\n')
						if err != nil {
							if err == io.EOF {
								// End of file; wait for more data
								time.Sleep(100 * time.Millisecond)
								continue
							}
							log.Printf("Error reading file: %v", err)
							return
						}
						logChan <- line
					}
				}
			}
		}
	}()

	// Wait for done signal or context cancellation
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
