package shared

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// Logger untuk logging
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info mencetak pesan informasi ke stdout
func Info(msg string) {
	InfoLogger.Println(msg)
}

// Error mencetak pesan kesalahan ke stderr
func Error(msg string) {
	ErrorLogger.Println(msg)
}

// RunCommand menjalankan perintah dan menangkap outputnya
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	// Menggunakan pipe untuk membaca stdout dan stderr secara bersamaan
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("command failed to start: %w", err)
	}

	// Membaca stdout dan stderr secara paralel
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go func() {
		defer close(stdoutDone)
		_, _ = io.Copy(&stdoutBuf, stdoutPipe)
	}()

	go func() {
		defer close(stderrDone)
		_, _ = io.Copy(&stderrBuf, stderrPipe)
	}()

	// Tunggu hingga perintah selesai
	if err := cmd.Wait(); err != nil {
		// Menangani kesalahan perintah dan log output stderr jika ada
		if stderrBuf.Len() > 0 {
			Error(fmt.Sprintf("command stderr: %s", stderrBuf.String()))
		}
		return "", fmt.Errorf("command failed: %w", err)
	}

	// Tunggu hingga semua output terbaca
	<-stdoutDone
	<-stderrDone

	if stdoutBuf.Len() > 0 {
		Info(fmt.Sprintf("command stdout: %s", stdoutBuf.String()))
	}

	return stdoutBuf.String(), nil
}

func ExitWithError(message string, err error) {
	Error(fmt.Sprintf("%s: %v\n", message, err))
	os.Exit(1)
}