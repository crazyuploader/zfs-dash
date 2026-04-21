package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourname/zfs-dash/internal/config"
	"github.com/yourname/zfs-dash/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the ZFS dashboard web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		if len(cfg.Endpoints) == 0 {
			fmt.Fprintln(os.Stderr, "No endpoints configured. Use --endpoints or zfs-dash.yaml")
			os.Exit(1)
		}
		return server.Start(cfg)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
