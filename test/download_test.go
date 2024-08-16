package test

import (
	"os"
	"sbx/internal"
	"testing"
)

func TestDownloadFile_GithubRelease(t *testing.T) {
	// URL dari file yang akan diunduh
	versiCaddy, err := internal.GetLatestRelease("caddyserver", "caddy", "stable")
	if err != nil {
			t.Logf("Gagal mendapatkan versi terbaru: %v\n", err)
			return
		}
	url, err := internal.BuildDownloadURL("caddyserver", "caddy",versiCaddy)
	t.Log(url)
	// urlCaddy := "https://github.com/caddyserver/caddy/releases/download/v2.8.4/caddy_2.8.4_linux_amd64.tar.gz"
	// Lokasi file sementara untuk menyimpan hasil unduhan
	// filepath := "/home/5678cxz/test-folder/sing-box-1.9.3-linux-amd64.tar.gz"
	filepath := "/home/5678cxz/test-folder/caddy_2.8.4_linux_amd64.tar.gz"
	// Pastikan file dihapus sebelum memulai tes
	os.Remove(filepath)

	// Mengunduh file
	err2 := internal.DownloadFile(url, filepath)
	if err2 != nil {
		t.Fatalf("Download failed: %v", err2)
	}
	// err3 := download.DownloadFile(urlCaddy, fp)
	// if err3 != nil {
	// 	t.Fatalf("Download failed: %v", err3)
	// }

	// Memeriksa apakah file berhasil diunduh
	fileInfo, err2 := os.Stat(filepath)
	if err2 != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Memeriksa apakah file tidak kosong
	if fileInfo.Size() == 0 {
		t.Fatalf("Downloaded file is empty")
	}

	// Cleanup: Menghapus file setelah pengujian selesai
	// os.Remove(filepath)
}
