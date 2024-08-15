package singbox

import "github.com/spf13/cobra"

var AccountCmd = &cobra.Command{
	Use:   "user",
	Short: "Kelola user sing-box (add, remove, list)",
}