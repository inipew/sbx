package caddy

import (
	"fmt"
	"os"
	"runtime"
	"sbx/internal"

	"github.com/spf13/cobra"
)


var InstallCmd = &cobra.Command{
    Use:   "install",
    Short: "Unduh dan instal Caddy",
    Run: func(cmd *cobra.Command, args []string) {
        err := InstallCaddy()
        if err != nil {
            fmt.Println("Gagal menginstal Caddy:", err)
            return
        }
        fmt.Println("Caddy berhasil diinstal.")
    },
}

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
	if err := internal.DownloadFile(url, tmpFilePath); err != nil {
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