package caddy

import (
	"fmt"
	"sbx/pkg/caddyfile"

	"github.com/spf13/cobra"
)

// Konstanta untuk pesan kesalahan
const (
	errApplyCaddyfile   = "Gagal menerapkan Caddyfile: %v"
	successApplyMessage = "Caddyfile berhasil diterapkan."
	errReplaceDomain    = "Gagal mengganti domain di Caddyfile: %v"
	successReplaceMessage = "Domain berhasil diganti."
)

// Fungsi untuk menerapkan Caddyfile
func applyCaddyfile(cmd *cobra.Command, args []string) {
	if err := caddyfile.ApplyCaddyfile(); err != nil {
		fmt.Printf(errApplyCaddyfile, err)
		return
	}
	fmt.Println(successApplyMessage)
}

// Fungsi untuk mengganti domain di Caddyfile
func replaceDomain(cmd *cobra.Command, args []string) {
	oldDomain := args[0]
	newDomain := args[1]
	if err := caddyfile.ReplaceDomain(oldDomain, newDomain); err != nil {
		fmt.Printf(errReplaceDomain, err)
		return
	}
	fmt.Println(successReplaceMessage)
}

// Mendefinisikan perintah `caddyfile`
var CaddyfileCmd = &cobra.Command{
	Use:   "caddyfile",
	Short: "Kelola Caddyfile",
	Long:  "Perintah untuk mengelola Caddyfile, termasuk menerapkan dan mengganti domain.",
}

func init() {
	// Menambahkan perintah `apply` ke `caddyfileCmd`
	CaddyfileCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Terapkan Caddyfile",
		Run:   applyCaddyfile,
	})

	// Menambahkan perintah `domain` ke `caddyfileCmd`
	CaddyfileCmd.AddCommand(&cobra.Command{
		Use:   "domain [oldDomain] [newDomain]",
		Short: "Ganti domain di Caddyfile",
		Args:  cobra.ExactArgs(2),
		Run:   replaceDomain,
	})
}
