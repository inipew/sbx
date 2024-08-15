package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"golang.org/x/time/rate"
)

func DownloadFile(url, filepath string) error {
	// Konfigurasi timeout dan rate limiter
	client := &http.Client{
		Timeout: 10 * time.Minute, // Batas waktu 10 menit
	}
	limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 request per detik

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ { // Retry hingga 3 kali
		lastErr = downloadWithRetries(client, limiter, url, filepath)
		if lastErr == nil {
			return nil
		}
		fmt.Printf("Retrying download (%d/3): %v\n", attempt+1, lastErr)
		time.Sleep(2 * time.Second) // Tunggu sebelum retry
	}

	return fmt.Errorf("failed to download file after multiple attempts: %w", lastErr)
}

func downloadWithRetries(client *http.Client, limiter *rate.Limiter, url, filepath string) error {
	// Mengirimkan request GET ke URL
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	// Memeriksa status response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %s", resp.Status)
	}

	// Membuka file untuk ditulis
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	// Menyalin konten dari response ke file dengan progress indicator
	return copyWithProgress(out, resp.Body, limiter)
}

func copyWithProgress(dst io.Writer, src io.Reader, limiter *rate.Limiter) error {
	buf := make([]byte, 32*1024) // Buffer 32 KB
	totalBytes := int64(0)
	ctx := context.Background() // Menggunakan konteks default

	for {
		n, err := src.Read(buf)
		if n > 0 {
			// Mengatur kecepatan rate limiter
			if err := limiter.Wait(ctx); err != nil {
				return fmt.Errorf("rate limiter error: %w", err)
			}

			if _, err := dst.Write(buf[:n]); err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}

			totalBytes += int64(n)
			fmt.Printf("\rDownloaded %d bytes", totalBytes) // Progress indicator
		}

		if err == io.EOF {
			fmt.Println() // New line after progress indicator
			break
		}
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
	}

	return nil
}

// InstallCaddy mengunduh dan menginstal Caddy sesuai dengan sistem operasi
func InstallCaddy() error {
	var url string
	switch runtime.GOOS {
	case "linux":
		url = "https://caddyserver.com/api/download?os=linux&arch=amd64"
	case "darwin":
		url = "https://caddyserver.com/api/download?os=darwin&arch=amd64"
	default:
		return fmt.Errorf("sistem operasi tidak didukung")
	}

	tmpFilePath := "/tmp/caddy"
	if err := DownloadFile(url, tmpFilePath); err != nil {
		return err
	}

	if err := os.Chmod(tmpFilePath, 0755); err != nil {
		return fmt.Errorf("gagal mengatur izin file: %w", err)
	}

	installPath := "/usr/local/bin/caddy"
	if err := os.Rename(tmpFilePath, installPath); err != nil {
		return fmt.Errorf("gagal memindahkan file ke lokasi instalasi: %w", err)
	}

	fmt.Println("Caddy berhasil diinstal di", installPath)
	return nil
}
