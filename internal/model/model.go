// Package model defines ZFS domain types and metric extraction logic.
package model

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/parser"
)

// PoolHealth is the human-readable pool state.
type PoolHealth string

const (
	HealthOnline    PoolHealth = "ONLINE"
	HealthDegraded  PoolHealth = "DEGRADED"
	HealthFaulted   PoolHealth = "FAULTED"
	HealthOffline   PoolHealth = "OFFLINE"
	HealthUnavail   PoolHealth = "UNAVAIL"
	HealthRemoved   PoolHealth = "REMOVED"
	HealthSuspended PoolHealth = "SUSPENDED"
	HealthUnknown   PoolHealth = "UNKNOWN"
)

// Pool holds the key metrics for one ZFS pool.
type Pool struct {
	Name               string     `json:"name"`
	Health             PoolHealth `json:"health"`
	Size               float64    `json:"size"`
	Allocated          float64    `json:"allocated"`
	Free               float64    `json:"free"`
	Freeing            float64    `json:"freeing"`
	LeakedBytes        float64    `json:"leaked_bytes"`
	DedupRatio         float64    `json:"dedup_ratio"`
	FragmentationRatio float64    `json:"fragmentation_ratio"`
	ReadOnly           bool       `json:"read_only"`
	UsedPercent        float64    `json:"used_percent"`
	Datasets           []Dataset  `json:"datasets,omitempty"`
}

// Dataset holds filesystem, volume, or snapshot metrics within a pool.
type Dataset struct {
	Name              string  `json:"name"`
	Pool              string  `json:"pool"`
	Type              string  `json:"type"`
	Available         float64 `json:"available"`
	LogicalUsed       float64 `json:"logical_used"`
	Quota             float64 `json:"quota"`
	Referenced        float64 `json:"referenced"`
	UsedByDataset     float64 `json:"used_by_dataset"`
	Used              float64 `json:"used"`
	VolumeSize        float64 `json:"volume_size"`
	Written           float64 `json:"written"`
	UsedPercent       float64 `json:"used_percent"`
	QuotaUsedPercent  float64 `json:"quota_used_percent"`
	VolumeUsedPercent float64 `json:"volume_used_percent"`
}

// DiskInfo holds key SMART metrics for one physical disk.
type DiskInfo struct {
	Device          string  `json:"device"`
	ModelName       string  `json:"model_name,omitempty"`
	SerialNumber    string  `json:"serial_number,omitempty"`
	FirmwareVersion string  `json:"firmware_version,omitempty"`
	Interface       string  `json:"interface,omitempty"`   // "sat", "nvme", "scsi"
	FormFactor      string  `json:"form_factor,omitempty"` // "3.5 inches", ""
	Temperature     float64 `json:"temperature"`           // °C current, 0 if unknown
	TempMin         float64 `json:"temp_min,omitempty"`    // lifetime min °C
	TempMax         float64 `json:"temp_max,omitempty"`    // lifetime max °C
	TempTrip        float64 `json:"temp_trip,omitempty"`   // hardware trip/critical °C
	SmartPassed     bool    `json:"smart_passed"`
	PowerOnHours    float64 `json:"power_on_hours"`
	PowerCycles     float64 `json:"power_cycles,omitempty"`
	CapacityBytes   float64 `json:"capacity_bytes"`
	RotationRate    int     `json:"rotation_rate"`             // RPM; 0 = SSD/NVMe
	PercentageUsed  float64 `json:"percentage_used,omitempty"` // SSD/NVMe wear 0–100
	AvailableSpare  float64 `json:"available_spare,omitempty"` // NVMe spare %
	SpareThreshold  float64 `json:"spare_threshold,omitempty"` // NVMe spare threshold %
	MediaErrors      float64 `json:"media_errors,omitempty"`     // unrecovered integrity errors
	HasMediaErrors   bool    `json:"has_media_errors,omitempty"` // true when metric was present (even if 0)
	CriticalWarning  float64 `json:"critical_warning,omitempty"`
	BytesRead        float64 `json:"bytes_read,omitempty"`
	BytesWritten     float64 `json:"bytes_written,omitempty"`
	ErrorLogCount        float64 `json:"error_log_count,omitempty"`
	HasPercentUsed       bool    `json:"has_percent_used,omitempty"`  // true when metric was present (even if 0)
	ReallocatedSectors   float64 `json:"reallocated_sectors,omitempty"`
	PendingSectors       float64 `json:"pending_sectors,omitempty"`
	OfflineUncorrectable float64 `json:"offline_uncorrectable,omitempty"`
	ReportedUncorrect    float64 `json:"reported_uncorrect,omitempty"`
	UDMACRCErrors        float64 `json:"udma_crc_errors,omitempty"`
	LoadCycleCount       float64 `json:"load_cycle_count,omitempty"`
	InterfaceSpeed       float64 `json:"interface_speed,omitempty"` // bits/sec
	ExitStatus           float64 `json:"exit_status,omitempty"`
	HasExitStatus        bool    `json:"has_exit_status,omitempty"`
}

// ExporterInfo holds metadata from zfs_exporter_build_info.
type ExporterInfo struct {
	Version   string `json:"version,omitempty"`
	GoVersion string `json:"go_version,omitempty"`
	Revision  string `json:"revision,omitempty"`
}

// NodeData holds all pool data fetched from one endpoint.
type NodeData struct {
	Label        string       `json:"label"`
	Location     string       `json:"location,omitempty"`
	URL          string       `json:"url"`
	FetchedAt    time.Time    `json:"fetched_at"`
	Error        string       `json:"error,omitempty"`
	ExporterInfo ExporterInfo `json:"exporter_info,omitempty"`
	Pools        []Pool       `json:"pools"`
	Disks        []DiskInfo   `json:"disks,omitempty"`
}

func healthFromValue(v float64) PoolHealth {
	switch int(math.Round(v)) {
	case 0:
		return HealthOnline
	case 1:
		return HealthDegraded
	case 2:
		return HealthFaulted
	case 3:
		return HealthOffline
	case 4:
		return HealthUnavail
	case 5:
		return HealthRemoved
	case 6:
		return HealthSuspended
	default:
		return HealthUnknown
	}
}

// ExtractExporterInfo extracts metadata from the zfs_exporter_build_info sample.
func ExtractExporterInfo(samples []parser.Sample) ExporterInfo {
	for _, s := range samples {
		if s.Name == "zfs_exporter_build_info" {
			return ExporterInfo{
				Version:   s.Labels["version"],
				GoVersion: s.Labels["goversion"],
				Revision:  s.Labels["revision"],
			}
		}
	}
	return ExporterInfo{}
}

// ExtractDisks builds DiskInfo structs from smartctl_exporter Prometheus samples.
// The mapping between pool and disk is not available in the metrics, so disks
// are returned at node level only.
func ExtractDisks(samples []parser.Sample) []DiskInfo {
	disks := map[string]*DiskInfo{}

	ensure := func(device string) *DiskInfo {
		if _, ok := disks[device]; !ok {
			disks[device] = &DiskInfo{Device: device, SmartPassed: true}
		}
		return disks[device]
	}

	for _, s := range samples {
		device := s.Labels["device"]
		if device == "" {
			continue
		}
		switch s.Name {
		case "smartctl_device":
			d := ensure(device)
			if d.ModelName == "" {
				d.ModelName = s.Labels["model_name"]
			}
			if d.SerialNumber == "" {
				d.SerialNumber = s.Labels["serial_number"]
			}
			if d.FirmwareVersion == "" {
				d.FirmwareVersion = s.Labels["firmware_version"]
			}
			if d.Interface == "" {
				d.Interface = s.Labels["interface"]
			}
			if d.FormFactor == "" {
				d.FormFactor = s.Labels["form_factor"]
			}
		case "smartctl_device_temperature":
			d := ensure(device)
			switch s.Labels["temperature_type"] {
			case "current":
				d.Temperature = s.Value
			case "drive_trip":
				d.TempTrip = s.Value
			case "lifetime_min":
				if d.TempMin == 0 || s.Value < d.TempMin {
					d.TempMin = s.Value
				}
			case "lifetime_max":
				if d.TempMax == 0 || s.Value > d.TempMax {
					d.TempMax = s.Value
				}
			}
		case "smartctl_device_smart_status":
			ensure(device).SmartPassed = s.Value == 1
		case "smartctl_device_power_on_seconds":
			ensure(device).PowerOnHours = s.Value / 3600
		case "smartctl_device_power_cycle_count":
			ensure(device).PowerCycles = s.Value
		case "smartctl_device_capacity_bytes":
			ensure(device).CapacityBytes = s.Value
		case "smartctl_device_rotation_rate":
			ensure(device).RotationRate = int(math.Round(s.Value))
		case "smartctl_device_percentage_used":
			d := ensure(device)
			d.PercentageUsed = s.Value
			d.HasPercentUsed = true
		case "smartctl_device_available_spare":
			ensure(device).AvailableSpare = s.Value
		case "smartctl_device_available_spare_threshold":
			ensure(device).SpareThreshold = s.Value
		case "smartctl_device_media_errors":
			d := ensure(device)
			d.MediaErrors = s.Value
			d.HasMediaErrors = true
		case "smartctl_device_critical_warning":
			ensure(device).CriticalWarning = s.Value
		case "smartctl_device_bytes_read":
			ensure(device).BytesRead = s.Value
		case "smartctl_device_bytes_written":
			ensure(device).BytesWritten = s.Value
		case "smartctl_device_error_log_count":
			if s.Labels["error_log_type"] == "summary" {
				ensure(device).ErrorLogCount = s.Value
			}
		case "smartctl_device_num_err_log_entries":
			// NVMe error log entries (separate metric from HDD error_log_count)
			d := ensure(device)
			if d.ErrorLogCount == 0 {
				d.ErrorLogCount = s.Value
			}
		case "smartctl_device_interface_speed":
			if s.Labels["speed_type"] == "current" {
				ensure(device).InterfaceSpeed = s.Value
			}
		case "smartctl_device_smartctl_exit_status":
			d := ensure(device)
			d.ExitStatus = s.Value
			d.HasExitStatus = true
		case "smartctl_device_attribute":
			if s.Labels["attribute_value_type"] != "raw" {
				continue
			}
			d := ensure(device)
			switch s.Labels["attribute_name"] {
			case "Reallocated_Sector_Ct":
				d.ReallocatedSectors = s.Value
			case "Current_Pending_Sector":
				d.PendingSectors = s.Value
			case "Offline_Uncorrectable":
				d.OfflineUncorrectable = s.Value
			case "Reported_Uncorrect":
				d.ReportedUncorrect = s.Value
			case "UDMA_CRC_Error_Count":
				d.UDMACRCErrors = s.Value
			case "Load_Cycle_Count":
				d.LoadCycleCount = s.Value
			}
		}
	}

	if len(disks) == 0 {
		return nil
	}

	result := make([]DiskInfo, 0, len(disks))
	for _, d := range disks {
		result = append(result, *d)
	}
	slices.SortFunc(result, func(a, b DiskInfo) int {
		return strings.Compare(a.Device, b.Device)
	})
	return result
}

// ExtractPools builds Pool structs from a flat Prometheus sample slice.
func ExtractPools(samples []parser.Sample) []Pool {
	pools := map[string]*Pool{}
	datasets := map[string]*Dataset{}
	ensure := func(name string) *Pool {
		if _, ok := pools[name]; !ok {
			pools[name] = &Pool{Name: name, Health: HealthUnknown}
		}
		return pools[name]
	}
	ensureDataset := func(pool, name, datasetType string) *Dataset {
		key := pool + "/" + name
		if _, ok := datasets[key]; !ok {
			datasets[key] = &Dataset{
				Name: name,
				Pool: pool,
				Type: datasetType,
			}
		}
		if datasetType != "" {
			datasets[key].Type = datasetType
		}
		return datasets[key]
	}

	for _, s := range samples {
		pool := s.Labels["pool"]
		if pool == "" {
			continue
		}
		p := ensure(pool)
		if datasetName := s.Labels["name"]; datasetName != "" {
			d := ensureDataset(pool, datasetName, s.Labels["type"])
			switch s.Name {
			case "zfs_dataset_available_bytes":
				d.Available = s.Value
			case "zfs_dataset_logical_used_bytes":
				d.LogicalUsed = s.Value
			case "zfs_dataset_quota_bytes":
				d.Quota = s.Value
			case "zfs_dataset_referenced_bytes":
				d.Referenced = s.Value
			case "zfs_dataset_used_by_dataset_bytes":
				d.UsedByDataset = s.Value
			case "zfs_dataset_used_bytes":
				d.Used = s.Value
			case "zfs_dataset_volume_size_bytes":
				d.VolumeSize = s.Value
			case "zfs_dataset_written_bytes":
				d.Written = s.Value
			}
		}

		switch s.Name {
		case "zfs_pool_health":
			p.Health = healthFromValue(s.Value)
		case "zfs_pool_size_bytes":
			p.Size = s.Value
		case "zfs_pool_allocated_bytes":
			p.Allocated = s.Value
		case "zfs_pool_free_bytes":
			p.Free = s.Value
		case "zfs_pool_freeing_bytes":
			p.Freeing = s.Value
		case "zfs_pool_leaked_bytes":
			p.LeakedBytes = s.Value
		case "zfs_pool_deduplication_ratio":
			p.DedupRatio = s.Value
		case "zfs_pool_fragmentation_ratio":
			p.FragmentationRatio = s.Value
		case "zfs_pool_readonly":
			p.ReadOnly = s.Value != 0
		}
	}

	result := make([]Pool, 0, len(pools))
	for _, d := range datasets {
		if d.Used > 0 && d.Available > 0 {
			total := d.Used + d.Available
			d.UsedPercent = math.Round((d.Used/total)*10000) / 100
		}
		if d.Quota > 0 {
			d.QuotaUsedPercent = math.Round((d.Used/d.Quota)*10000) / 100
		}
		if d.VolumeSize > 0 {
			d.VolumeUsedPercent = math.Round((d.Used/d.VolumeSize)*10000) / 100
		}
		if p, ok := pools[d.Pool]; ok {
			p.Datasets = append(p.Datasets, *d)
		}
	}

	for _, p := range pools {
		if p.Size > 0 {
			if p.Allocated > 0 && p.Free == 0 {
				p.Free = p.Size - p.Allocated
			} else if p.Free > 0 && p.Allocated == 0 {
				p.Allocated = p.Size - p.Free
			}
			p.UsedPercent = math.Round((p.Allocated/p.Size)*10000) / 100
		}
		slices.SortFunc(p.Datasets, func(a, b Dataset) int {
			if a.Used != b.Used {
				if a.Used > b.Used {
					return -1
				}
				return 1
			}
			return strings.Compare(a.Name, b.Name)
		})
		result = append(result, *p)
	}

	slices.SortFunc(result, func(a, b Pool) int {
		return strings.Compare(a.Name, b.Name)
	})

	return result
}

// HumanBytes returns a human-readable byte string (e.g. "3.72 TB").
func HumanBytes(b float64) string {
	const unit = 1024.0
	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	if b < unit {
		return fmt.Sprintf("%.0f B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit && exp < len(units)-1; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %s", b/div, units[exp])
}
