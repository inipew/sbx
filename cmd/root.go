package cmd

import (
	"sbx/cmd/caddy"
	"sbx/cmd/install"
	"sbx/cmd/singbox"

	"github.com/spf13/cobra"
)
func NewRootCmd() *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "manager",
        Short: "Manajemen aplikasi CLI",
        Long:  "Aplikasi CLI untuk manajemen Caddy dan Sing-box.",
    }

    // Tambahkan subcommand untuk Caddy
    caddyCmd := &cobra.Command{
        Use:   "caddy",
        Short: "Manajemen Caddy",
        Long:  "Subcommand untuk manajemen Caddy termasuk update, install, service, dan konfigurasi Caddyfile.",
    }
    caddyCmd.AddCommand(caddy.UpdateCmd)
    caddyCmd.AddCommand(caddy.InstallCmd)
    caddyCmd.AddCommand(caddy.ServiceCmd)
    caddyCmd.AddCommand(caddy.CaddyfileCmd)
    caddyCmd.AddCommand(caddy.LogCmd)
    caddyCmd.AddCommand(caddy.CheckUpdateCmd)

    // Tambahkan subcommand untuk Sing
    singCmd := &cobra.Command{
        Use:   "sing",
        Short: "Manajemen Sing",
        Long:  "Subcommand untuk manajemen Sing-box termasuk update, install, service, dan manajemen account.",
    }
    singCmd.AddCommand(singbox.InstallCmd)
    singCmd.AddCommand(singbox.UpdateCmd)
    singCmd.AddCommand(singbox.ServiceCmd)
    singCmd.AddCommand(singbox.AccountCmd)
    singCmd.AddCommand(singbox.LogCmd)
    singCmd.AddCommand(singbox.CheckUpdateCmd)

    installCmd := &cobra.Command{
        Use:   "install",
        Short: "Install Service",
        Long:  "Subcommand untuk menginstall all service.",
    }
    installCmd.AddCommand(install.InstallAllCmd)
    installCmd.AddCommand(install.InstallCaddyCmd)
    installCmd.AddCommand(install.InstallSingCmd)

    rootCmd.AddCommand(caddyCmd)
    rootCmd.AddCommand(singCmd)
    rootCmd.AddCommand(installCmd)

    return rootCmd
}
