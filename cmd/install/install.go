package install

import (
	"fmt"
	"os"
	"sbx/internal"
	"sbx/shared"

	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
    Use:   "install [all]",
    Short: "Install Service",
}

var InstallAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Install semua layanan",
	Run:   func(cmd *cobra.Command, args []string) {
		err := internal.RestartService("sing-box")
		if err != nil {
			shared.Error(fmt.Sprintf("Error watching log file: %v", err))
		}
	},
}

var InstallCaddyCmd = &cobra.Command{
	Use:   "caddy",
	Short: "Install layanan caddy",
	Run:   func(cmd *cobra.Command, args []string) {
		err := internal.RestartService("sing-box")
		if err != nil {
			shared.Error(fmt.Sprintf("Error watching log file: %v", err))
		}
	},
}

var InstallSingCmd = &cobra.Command{
	Use:   "sing",
	Short: "Install layanan sing-box",
	Run:   func(cmd *cobra.Command, args []string) {
		err := internal.RestartService("sing-box")
		if err != nil {
			shared.Error(fmt.Sprintf("Error watching log file: %v", err))
		}
	},
}

func create(cmd *cobra.Command, args []string){
	dirs := []string{
		internal.BinDir,
		internal.BackupDir,
		internal.TmpDir,
		internal.SingboxConfDir,
		internal.CaddyDir,
	}

	// Panggil fungsi untuk membuat semua folder
	if err := createAllFolders(dirs); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// createAllFolders membuat folder dari daftar dan menetapkan izin 755
func createAllFolders(dirs []string) error {
	for _, dir := range dirs {
		if err := createDirWithPermission(dir); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
		fmt.Printf("Directory created and permission set: %s\n", dir)
	}
	return nil
}

// createDirWithPermission membuat folder dan menetapkan izin 755
func createDirWithPermission(dir string) error {
	// Buat folder dan semua folder yang diperlukan
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Set izin folder ke 755 (rwxr-xr-x)
	err = os.Chmod(dir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func CreateDomainFile(domain string) error {
	err := os.WriteFile(internal.DomainFilePath, []byte(domain), 0644)
	if err != nil {
		return fmt.Errorf("failed to write domain file: %w", err)
	}
	return nil
}

func installCaddy(){

}

func CreateAndEnableService() error {
	err := os.WriteFile(internal.CaddyServicePath, []byte(internal.CaddyServiceContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}

	return internal.RunSystemdCommand("caddy.service", "enable")
}