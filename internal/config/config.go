package config

import (
	"cmp"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Endpoint is a single ZFS exporter target with optional companion exporters.
type Endpoint struct {
	URL             string `mapstructure:"url"`
	Label           string `mapstructure:"label"`
	Location        string `mapstructure:"location"`
	SmartctlURL     string `mapstructure:"smartctl_url"`      // optional; omit to skip disk metrics
	NodeExporterURL string `mapstructure:"node_exporter_url"` // optional; omit to skip system metrics
}

// HistoryConfig controls the embedded time-series store.
type HistoryConfig struct {
	Enabled        bool
	Path           string
	Retention      time.Duration
	RecordInterval time.Duration
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
	History         HistoryConfig
}

// parseFlexDuration converts a config value into a duration.
// Accepts bare numbers (seconds, back-compat), numeric strings (seconds),
// and Go duration strings ("5m", "1h30m"). Returns def when v is nil,
// unparseable, or non-positive.
func parseFlexDuration(v any, def time.Duration) time.Duration {
	var d time.Duration
	switch t := v.(type) {
	case nil:
		return def
	case int:
		d = time.Duration(t) * time.Second
	case int64:
		d = time.Duration(t) * time.Second
	case float64:
		d = time.Duration(t * float64(time.Second))
	case string:
		if secs, err := strconv.Atoi(t); err == nil {
			d = time.Duration(secs) * time.Second
		} else if parsed, err := time.ParseDuration(t); err == nil {
			d = parsed
		} else {
			return def
		}
	default:
		return def
	}
	if d <= 0 {
		return def
	}
	return d
}

// Load reads viper state into a validated Config.
func Load() (*Config, error) {
	histRetention := viper.GetDuration("history.retention")
	if histRetention <= 0 {
		histRetention = 720 * time.Hour // 30 days default
	}
	cfg := &Config{
		Addr:            viper.GetString("addr"),
		Refresh:         parseFlexDuration(viper.Get("refresh"), 300*time.Second),
		CacheTTL:        time.Duration(cmp.Or(viper.GetInt("cache_ttl"), 30)) * time.Second,
		Debug:           viper.GetBool("debug"),
		TrustedProxies:  viper.GetStringSlice("trusted_proxies"),
		MaxUsagePercent: viper.GetFloat64("max_usage_percent"),
		LogFormat:       cmp.Or(viper.GetString("log_format"), "text"),
		History: HistoryConfig{
			Enabled:        viper.GetBool("history.enabled"),
			Path:           cmp.Or(viper.GetString("history.path"), "./data/history.db"),
			Retention:      histRetention,
			RecordInterval: viper.GetDuration("history.record_interval"),
		},
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
