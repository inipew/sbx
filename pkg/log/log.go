package log

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// watchLog monitors the specified log file for changes and prints new log entries.
func WatchLog(path string) error {
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

    err = watcher.Add(path)
    if err != nil {
        return fmt.Errorf("failed to add file to watcher: %w", err)
    }

    // Use a buffered reader
    reader := bufio.NewReader(file)
    file.Seek(0, io.SeekEnd)

    logChan := make(chan string)

    go func() {
        for line := range logChan {
            fmt.Print(line)
        }
    }()

    go func() {
        for event := range watcher.Events {
            if event.Op&fsnotify.Write == fsnotify.Write {
                for {
                    line, err := reader.ReadString('\n')
                    if err != nil {
                        if err == io.EOF {
                            time.Sleep(100 * time.Millisecond) // Small sleep to reduce CPU usage
                            continue
                        }
                        log.Printf("Error reading file: %v", err)
                        break
                    }
                    logChan <- line
                }
            }
        }
    }()

    // Keep the program running to continue watching the file
    select {}
}