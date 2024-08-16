package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadFile(url, filepath string) error {
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		lastErr = downloadFile(client, url, filepath)
		if lastErr == nil {
			return nil
		}
		fmt.Printf("Retrying download (%d/3): %v\n", attempt+1, lastErr)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to download file after multiple attempts: %w", lastErr)
}

func downloadFile(client *http.Client, url, filepath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	return copyWithProgress(out, resp.Body)
}

func copyWithProgress(dst io.Writer, src io.Reader) error {
	buf := make([]byte, 64*1024) // Buffer 64 KB
	totalBytes := int64(0)
	lastPrintTime := time.Now()
	startTime := time.Now()
	bytesSinceLastPrint := int64(0)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}

			totalBytes += int64(n)
			bytesSinceLastPrint += int64(n)

			// Mengupdate progress setiap detik
			if time.Since(lastPrintTime) > time.Second {
				duration := time.Since(startTime).Seconds()
				speed := float64(bytesSinceLastPrint) / duration // kecepatan dalam byte per detik

				fmt.Printf("\rDownloaded %d bytes (%.2f MB), Speed: %.2f MB/s", totalBytes, float64(totalBytes)/1e6, speed/1e6)
				lastPrintTime = time.Now()
				bytesSinceLastPrint = 0
			}
		}

		if err == io.EOF {
			fmt.Println()
			break
		}
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
	}

	return nil
}
