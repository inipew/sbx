package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReplaceDomain(oldDomain, newDomain string) error {
	var oldDomainStr string

	if oldDomain == "" {
		// Read the old domain from the domain file if oldDomain is empty
		data, err := os.ReadFile(DomainFilePath)
		if err != nil {
			return fmt.Errorf("failed to read domain file: %w", err)
		}
		oldDomainStr = strings.TrimSpace(string(data))
	} else {
		// Use the provided old domain
		oldDomainStr = oldDomain
	}
    // Baca isi Caddyfile
    data, err := os.ReadFile(CaddyFilePath)
    if err != nil {
        return fmt.Errorf("gagal membaca Caddyfile: %w", err)
    }

    // Ganti domain lama dengan domain baru
    updatedData := strings.ReplaceAll(string(data), oldDomainStr, newDomain)

    // Tulis kembali ke Caddyfile
    err = os.WriteFile(CaddyFilePath, []byte(updatedData), 0644)
    if err != nil {
        return fmt.Errorf("gagal menulis Caddyfile: %w", err)
    }

    return ApplyCaddyfile()
}

// InstallDefaultConfig menulis konfigurasi default ke Caddyfile
func InstallDefaultConfig() error {
    err := os.WriteFile(CaddyFilePath, []byte(CaddyFileContent), 0644)
    if err != nil {
        return fmt.Errorf("gagal menulis Caddyfile default: %w", err)
    }
    return nil
}

func ApplyCaddyfile() error {
    cmd := exec.Command(CaddyBinPath, "reload", "--config", CaddyFilePath, "--force")
    out, err := cmd.CombinedOutput() // Mengambil output error juga
    if err != nil {
        return fmt.Errorf("gagal memuat ulang Caddyfile: %s", string(out))
    }
    return nil
}