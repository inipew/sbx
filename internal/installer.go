package internal

import "fmt"


func InstallCaddy(version string) error {
	version, err := GetLatestRelease("caddyserver","caddy","stable")
	if err != nil {
		fmt.Printf("Gagal mendapatkan versi terbaru: %v\n", err)

	}
	url,err := BuildDownloadURL("caddyserver","caddy",version)
	if err != nil {
		fmt.Printf("Gagal download: %v\n", err)

	}

	if err := DownloadFile(url, TmpDir); err != nil {
		return fmt.Errorf("failed to download Caddy: %w", err)
	}

	if err := ExtractTarGz(TmpDir, CaddyBinPath); err != nil {
		return fmt.Errorf("failed to extract Caddy tarball: %w", err)
	}

	if err := CreateServiceFile("caddy"); err != nil {
		return fmt.Errorf("failed to create Caddy service: %w", err)
	}

	if err := EnableService("caddy.service"); err != nil {
		return fmt.Errorf("failed to enable Caddy service: %w", err)
	}
	if err := StartService("caddy.service"); err != nil {
		return fmt.Errorf("failed to start Caddy service: %w", err)
	}

	if err := InstallDefaultConfig(); err != nil {
		return fmt.Errorf("failed to create Caddyfile: %w", err)
	}
	if err := ApplyCaddyfile(); err != nil {
		return fmt.Errorf("failed to apply Caddyfile: %w", err)
	}

	return nil
}