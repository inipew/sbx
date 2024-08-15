package main

import (
	"fmt"
	"os"
	"sbx/cmd"
)

func main() {
    // Buat command root
    rootCmd := cmd.NewRootCmd()

    // Jalankan command root
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1) // Keluar dengan status non-nol jika terjadi kesalahan
    }
}
