package caddy

import (
	"fmt"
	"sbx/pkg/download"

	"github.com/spf13/cobra"
)


var InstallCmd = &cobra.Command{
    Use:   "install",
    Short: "Unduh dan instal Caddy",
    Run: func(cmd *cobra.Command, args []string) {
        err := download.InstallCaddy()
        if err != nil {
            fmt.Println("Gagal menginstal Caddy:", err)
            return
        }
        fmt.Println("Caddy berhasil diinstal.")
    },
}
