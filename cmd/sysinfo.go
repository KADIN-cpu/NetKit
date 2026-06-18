package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KADIN-cpu/netkit/internal/output"
	"github.com/KADIN-cpu/netkit/internal/sysinfo"
)

var sysinfoCmd = &cobra.Command{
	Use:     "sysinfo",
	Aliases: []string{"sys", "info"},
	Short:   "Mostra informações do sistema (CPU, memória, disco, uptime)",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := sysinfo.Collect()
		if err != nil {
			return err
		}

		output.Section("Sistema")
		output.Table(
			[]string{"Propriedade", "Valor"},
			[][]string{
				{"Hostname", info.Hostname},
				{"SO", info.OS},
				{"Plataforma", info.Platform},
				{"Arquitetura", info.KernelArch},
				{"Uptime", sysinfo.HumanDuration(info.Uptime)},
			},
		)

		output.Section("CPU")
		output.Table(
			[]string{"Propriedade", "Valor"},
			[][]string{
				{"Modelo", info.CPUModel},
				{"Núcleos", fmt.Sprintf("%d", info.CPUCores)},
				{"Uso", fmt.Sprintf("%.1f%%", info.CPUPercent)},
			},
		)

		output.Section("Memória")
		output.Table(
			[]string{"Propriedade", "Valor"},
			[][]string{
				{"Total", sysinfo.HumanBytes(info.MemTotal)},
				{"Em uso", sysinfo.HumanBytes(info.MemUsed)},
				{"Uso", fmt.Sprintf("%.1f%%", info.MemPercent)},
			},
		)

		output.Section("Discos")
		rows := make([][]string, 0, len(info.Disks))
		for _, d := range info.Disks {
			rows = append(rows, []string{
				d.Mountpoint,
				d.Fstype,
				sysinfo.HumanBytes(d.Total),
				sysinfo.HumanBytes(d.Used),
				fmt.Sprintf("%.1f%%", d.UsedPercent),
			})
		}
		output.Table([]string{"Ponto de montagem", "FS", "Total", "Usado", "Uso"}, rows)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sysinfoCmd)
}
