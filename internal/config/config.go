package config

import (
	"cmp"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Endpoint is a single ZFS exporter target.
type Endpoint struct {
	URL      string `mapstructure:"url"`
	Label    string `mapstructure:"label"`
	Location string `mapstructure:"location"`
}

// Config holds all runtime options.
type Config struct {
	Endpoints       []Endpoint
	Addr            string
	Refresh         time.Duration
	CacheTTL        time.Duration
	Debug           bool
	TrustedProxies  []string
	MaxUsagePercent float64
	LogFormat       string // "text" or "json"
}

// Load reads viper state into a validated Config.
func Load() (*Config, error) {
	cfg := &Config{
		Addr:            cmp.Or(viper.GetString("addr"), ":8054"),
		Refresh:         time.Duration(viper.GetInt("refresh")) * time.Second,
		CacheTTL:        time.Duration(cmp.Or(viper.GetInt("cache_ttl"), 30)) * time.Second,
		Debug:           viper.GetBool("debug"),
		TrustedProxies:  viper.GetStringSlice("trusted_proxies"),
		MaxUsagePercent: viper.GetFloat64("max_usage_percent"),
		LogFormat:       cmp.Or(viper.GetString("log_format"), "text"),
	}
	if cfg.Refresh <= 0 {
		cfg.Refresh = 300 * time.Second
	}

	// Try structured endpoints block (config file).
	var eps []Endpoint
	if viper.IsSet("endpoints") {
		if err := viper.UnmarshalKey("endpoints", &eps); err != nil {
			return nil, fmt.Errorf("decode endpoints: %w", err)
		}

		if len(eps) > 0 {
			for i, ep := range eps {
				if ep.URL == "" {
					return nil, fmt.Errorf("endpoint[%d] missing url", i)
				}
				if ep.Label == "" {
					eps[i].Label = ep.URL
				}
			}
			cfg.Endpoints = eps
			return cfg, nil
		}
	}

	// Fall back to flat string slice (--endpoints flag / env).
	for _, u := range viper.GetStringSlice("endpoints") {
		if u != "" {
			cfg.Endpoints = append(cfg.Endpoints, Endpoint{URL: u, Label: u})
		}
	}
	return cfg, nil
}
