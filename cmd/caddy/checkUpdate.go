package caddy

import (
	"fmt"
	"sbx/internal"
	"sbx/shared"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
)

// Mendefinisikan perintah `checkUpdate`
var CheckUpdateCmd = &cobra.Command{
	Use:   "checkupdate [type]",
	Short: "Periksa apakah ada pembaruan aplikasi",
	Long:  "Perintah untuk memeriksa apakah ada versi terbaru dari aplikasi yang tersedia.",
	Args:  cobra.ExactArgs(1),
	Run:   checkUpdate,
}

// Fungsi untuk memeriksa pembaruan
func checkUpdate(cmd *cobra.Command, args []string) {
	jenisRilis := args[0]

	if jenisRilis != "stable" && jenisRilis != "latest" {
		fmt.Println("Jenis rilis tidak valid. Pilih 'stable' atau 'latest'.")
		return
	}

	// Mendapatkan versi terbaru dari Sing-box
	latestVersion, err := internal.GetLatestRelease("caddyserver", "caddy", jenisRilis)
	if err != nil {
		fmt.Printf("Gagal mendapatkan versi terbaru: %v\n", err)
		return
	}

	// Mendapatkan versi saat ini dari Sing-box
	currentVersion, err := internal.GetCaddyVersion()
	if err != nil {
		shared.Info(fmt.Sprintf("Gagal mendapatkan versi saat ini: %v\n", err))
		return
	}

	// Parse versi untuk membandingkan
	latestSemver, err := semver.Parse(latestVersion)
	if err != nil {
		shared.Info(fmt.Sprintf("Gagal parsing versi terbaru: %v\n", err))
		return
	}

	currentSemver, err := semver.Parse(currentVersion)
	if err != nil {
		shared.Info(fmt.Sprintf("Gagal parsing versi saat ini: %v\n", err))
		return
	}

	// Bandingkan versi
	if latestSemver.GT(currentSemver) {
		shared.Info(fmt.Sprintln("Ada pembaruan tersedia!"))
		shared.Info(fmt.Sprintf("Versi saat ini: %s | Versi terbaru: %s", currentVersion, latestVersion))
	} else {
		shared.Info(fmt.Sprintf("Versi saat ini: %s | Versi terbaru: %s", currentVersion, latestVersion))
		shared.Info(fmt.Sprintln("Aplikasi Anda sudah versi terbaru."))
	}
}
