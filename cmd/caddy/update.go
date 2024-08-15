package caddy

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sbx/pkg/github"
	service "sbx/pkg/service/caddy"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Periksa dan perbarui Caddy ke versi terbaru",
	Run: func(cmd *cobra.Command, args []string) {
		// Mendapatkan versi terbaru dari Caddy
		// Menentukan jenis rilis dari argumen atau default ke "stable"
		jenisRilis, _ := cmd.Flags().GetString("release-type")
		if jenisRilis == "" {
			jenisRilis = "stable"
		}

		// Mendapatkan versi terbaru dari Sing-box
		latestVersion, err := github.GetLatestRelease("SagerNet", "sing-box", jenisRilis)
		if err != nil {
			fmt.Printf("Gagal mendapatkan versi terbaru: %v\n", err)
			return
		}
		fmt.Printf("Versi terbaru Caddy adalah: %s\n", latestVersion)

		// Mendapatkan versi saat ini dari Caddy
		currentVersion, err := service.GetCurrentVersion()
		if err != nil {
			fmt.Printf("Gagal mendapatkan versi saat ini: %v\n", err)
			return
		}
		fmt.Printf("Versi saat ini Caddy adalah: %s\n", currentVersion)

		// Membandingkan versi saat ini dengan versi terbaru
		if currentVersion == latestVersion {
			fmt.Println("Caddy sudah diperbarui ke versi terbaru.")
			return
		}

		// Jika versi terbaru lebih tinggi, lakukan pembaruan
		fmt.Println("Memperbarui Caddy ke versi terbaru...")
		err = updateCaddy(latestVersion)
		if err != nil {
			fmt.Printf("Gagal memperbarui Caddy: %v\n", err)
			return
		}
		fmt.Printf("Caddy berhasil diperbarui ke versi: %s\n", latestVersion)
	},
}

// Fungsi untuk memperbarui Caddy
func updateCaddy(version string) error {
	// Menentukan URL unduhan berdasarkan versi dan sistem operasi
	var url string
	switch runtime.GOOS {
	case "linux":
		url = fmt.Sprintf("https://caddyserver.com/api/download?os=linux&arch=amd64&version=%s", version)
	case "darwin":
		url = fmt.Sprintf("https://caddyserver.com/api/download?os=darwin&arch=amd64&version=%s", version)
	default:
		return fmt.Errorf("sistem operasi tidak didukung untuk pembaruan")
	}

	// Unduh binari Caddy terbaru
	fmt.Printf("Mengunduh Caddy versi %s dari %s...\n", version, url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("gagal mengunduh binari Caddy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal mengunduh Caddy: status %s", resp.Status)
	}

	tempFile, err := os.CreateTemp("", "caddy_update_")
	if err != nil {
		return fmt.Errorf("gagal membuat file sementara: %w", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("gagal menulis file sementara: %w", err)
	}

	// Atur izin file dan ganti binari yang lama
	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		return fmt.Errorf("gagal mengatur izin file: %w", err)
	}

	// Lokasi binari Caddy saat ini
	currentPath, err := exec.LookPath("caddy")
	if err != nil {
		return fmt.Errorf("gagal menemukan binari Caddy saat ini: %w", err)
	}

	backupPath := currentPath + ".bak"
	if err := os.Rename(currentPath, backupPath); err != nil {
		return fmt.Errorf("gagal membackup binari Caddy lama: %w", err)
	}

	if err := os.Rename(tempFile.Name(), currentPath); err != nil {
		// Kembalikan binari yang lama jika mengganti gagal
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("gagal mengganti binari Caddy: %w", err)
	}

	// Hapus backup jika semuanya berhasil
	defer os.Remove(backupPath)

	// Verifikasi update
	updatedVersion, err := service.GetCurrentVersion()
	if err != nil || updatedVersion != version {
		return fmt.Errorf("pembaruan gagal, versi saat ini tidak sesuai")
	}

	fmt.Println("Caddy berhasil diperbarui.")
	return nil
}
