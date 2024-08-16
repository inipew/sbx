package singbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sbx/internal"
	"sbx/shared"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update [type]",
	Short: "Periksa dan perbarui Sing-box ke versi terbaru",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jenisRilis := args[0]

		if jenisRilis != "stable" && jenisRilis != "latest" {
			shared.Info(fmt.Sprintln("Jenis rilis tidak valid. Pilih 'stable' atau 'latest'."))
			return
		}

		// Mendapatkan versi terbaru dari Sing-box
		latestVersion, err := internal.GetLatestRelease("SagerNet", "sing-box", jenisRilis)
		if err != nil {
			fmt.Printf("Gagal mendapatkan versi terbaru: %v\n", err)
			return
		}
		fmt.Printf("Versi terbaru Sing-box adalah: %s\n", latestVersion)

		// Mendapatkan versi saat ini dari Sing-box
		currentVersion, err := internal.GetSingBoxVersion()
		if err != nil {
			fmt.Printf("Gagal mendapatkan versi saat ini: %v\n", err)
			return
		}
		fmt.Printf("Versi saat ini Sing-box adalah: %s\n", currentVersion)

		// Membandingkan versi saat ini dengan versi terbaru
		if currentVersion == latestVersion {
			fmt.Println("Sing-box sudah diperbarui ke versi terbaru.")
			return
		}

		// Jika versi terbaru lebih tinggi, lakukan pembaruan
		fmt.Println("Memperbarui Sing-box ke versi terbaru...")
		err = updateSing(latestVersion)
		if err != nil {
			fmt.Printf("Gagal memperbarui Sing-box: %v\n", err)
			return
		}
		fmt.Printf("Sing-box berhasil diperbarui ke versi: %s\n", latestVersion)
	},
}

// Fungsi untuk memperbarui Sing-box
// UpdateSing memperbarui binari Sing-box ke versi terbaru
func updateSing(version string) error {
	// Menentukan URL unduhan
	url, err := internal.BuildDownloadURL("SagerNet", "sing-box", version)
	tempFilePath := filepath.Join(internal.TmpDir, "singbox_update.tar.gz")
	extractDir := filepath.Join(internal.TmpDir, "singbox_update")

	fmt.Printf("Mengunduh Sing-box versi %s dari %s...\n", version, url)
	if err := internal.DownloadFile(url, tempFilePath); err != nil {
		return fmt.Errorf("gagal mengunduh binari Sing-box: %w", err)
	}

	fmt.Printf("Mengekstrak file %s...\n", tempFilePath)
	if err := internal.ExtractTarGz(tempFilePath, extractDir); err != nil {
		return fmt.Errorf("gagal mengekstrak file: %w", err)
	}

	newBinaryPath := filepath.Join(extractDir, "sing-box")
	if _, err := os.Stat(newBinaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binari baru tidak ditemukan setelah ekstraksi")
	}

	currentPath, err := exec.LookPath("sing-box")
	if err != nil {
		return fmt.Errorf("gagal menemukan binari Sing-box saat ini: %w", err)
	}

	backupPath := currentPath + ".bak"
	if err := os.Rename(currentPath, backupPath); err != nil {
		return fmt.Errorf("gagal membackup binari Sing-box lama: %w", err)
	}

	if err := os.Rename(newBinaryPath, currentPath); err != nil {
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("gagal mengganti binari Sing-box: %w", err)
	}

	defer os.Remove(backupPath)
	defer os.RemoveAll(extractDir)
	defer os.Remove(tempFilePath)

	updatedVersion, err := internal.GetLatestRelease("SagerNet", "sing-box", "stable")
	if err != nil || updatedVersion != version {
		return fmt.Errorf("pembaruan gagal, versi saat ini tidak sesuai")
	}

	fmt.Println("Sing-box berhasil diperbarui.")
	return nil
}
