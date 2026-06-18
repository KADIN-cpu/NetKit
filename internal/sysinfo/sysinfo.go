package sysinfo

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// Info agrupa todos os dados coletados do sistema.
type Info struct {
	Hostname   string
	OS         string
	Platform   string
	KernelArch string
	Uptime     time.Duration
	CPUModel   string
	CPUCores   int
	CPUPercent float64
	MemTotal   uint64
	MemUsed    uint64
	MemPercent float64
	Disks      []DiskInfo
}

// DiskInfo representa o uso de uma partição.
type DiskInfo struct {
	Mountpoint  string
	Fstype      string
	Total       uint64
	Used        uint64
	UsedPercent float64
}

// Collect reúne as informações do host atual.
func Collect() (*Info, error) {
	hi, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("host info: %w", err)
	}

	info := &Info{
		Hostname:   hi.Hostname,
		OS:         hi.OS,
		Platform:   fmt.Sprintf("%s %s", hi.Platform, hi.PlatformVersion),
		KernelArch: hi.KernelArch,
		Uptime:     time.Duration(hi.Uptime) * time.Second,
	}

	if cpus, err := cpu.Info(); err == nil && len(cpus) > 0 {
		info.CPUModel = cpus[0].ModelName
	}
	if counts, err := cpu.Counts(true); err == nil {
		info.CPUCores = counts
	}
	// Amostra de 1s para um percentual de uso confiável.
	if pcts, err := cpu.Percent(time.Second, false); err == nil && len(pcts) > 0 {
		info.CPUPercent = pcts[0]
	}

	if vm, err := mem.VirtualMemory(); err == nil {
		info.MemTotal = vm.Total
		info.MemUsed = vm.Used
		info.MemPercent = vm.UsedPercent
	}

	if parts, err := disk.Partitions(false); err == nil {
		for _, p := range parts {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}
			// Ignora partições virtuais sem tamanho.
			if usage.Total == 0 {
				continue
			}
			info.Disks = append(info.Disks, DiskInfo{
				Mountpoint:  p.Mountpoint,
				Fstype:      p.Fstype,
				Total:       usage.Total,
				Used:        usage.Used,
				UsedPercent: usage.UsedPercent,
			})
		}
	}

	return info, nil
}

// HumanBytes formata bytes em unidade legível (KB, MB, GB...).
func HumanBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// HumanDuration formata uma duração em dias/horas/minutos.
func HumanDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}
