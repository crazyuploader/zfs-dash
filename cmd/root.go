package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var initConfigErr error

var rootCmd = &cobra.Command{
	Use:           "zfs-dash",
	Short:         "ZFS Dashboard — real-time pool monitoring",
	Long:          `Pull ZFS exporter metrics from multiple Prometheus endpoints and serve a minimal real-time dashboard.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ./config.yaml)")
	rootCmd.PersistentFlags().StringSlice("endpoints", nil, "ZFS exporter /metrics URLs (comma-separated or repeated)")
	rootCmd.PersistentFlags().String("addr", ":8080", "listen address")
	rootCmd.PersistentFlags().Int("refresh", 30, "auto-refresh interval in seconds")

	_ = viper.BindPFlag("endpoints", rootCmd.PersistentFlags().Lookup("endpoints"))
	_ = viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	_ = viper.BindPFlag("refresh", rootCmd.PersistentFlags().Lookup("refresh"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/zfs-dash")
	}
	viper.SetEnvPrefix("ZFSDASH")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config:", viper.ConfigFileUsed())
	} else if cfgFile != "" || !errors.As(err, new(viper.ConfigFileNotFoundError)) {
		initConfigErr = err
	}
}

func configInitError() error {
	return initConfigErr
}
