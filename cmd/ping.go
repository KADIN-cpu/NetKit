package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/KADIN-cpu/netkit/internal/netscan"
	"github.com/KADIN-cpu/netkit/internal/output"
	"github.com/KADIN-cpu/netkit/internal/ping"
)

var (
	pingTimeout int
	pingWorkers int
	pingOnlyUp  bool
)

var pingCmd = &cobra.Command{
	Use:   "ping <host | CIDR>",
	Short: "Faz ping em um host ou varre uma faixa de rede (CIDR)",
	Long: `Faz ping em um único host ou em toda uma faixa de rede usando CIDR.

Exemplos:
  netkit ping 192.168.0.10
  netkit ping 192.168.0.0/24
  netkit ping 192.168.0.0/24 --only-up`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		timeout := time.Duration(pingTimeout) * time.Millisecond

		var hosts []string
		// Detecta se é CIDR ou host único.
		if isCIDR(target) {
			h, err := netscan.HostsFromCIDR(target)
			if err != nil {
				return err
			}
			hosts = h
		} else {
			hosts = []string{target}
		}

		output.Section(fmt.Sprintf("Ping em %d host(s)", len(hosts)))
		results := ping.Sweep(hosts, timeout, pingWorkers)

		// Ordena: vivos primeiro.
		sort.SliceStable(results, func(i, j int) bool {
			if results[i].Alive != results[j].Alive {
				return results[i].Alive
			}
			return false
		})

		rows := make([][]string, 0, len(results))
		upCount := 0
		for _, r := range results {
			if r.Alive {
				upCount++
			} else if pingOnlyUp {
				continue
			}
			status := output.Down("offline")
			rtt := "-"
			if r.Alive {
				status = output.Up("online")
				rtt = fmt.Sprintf("%d ms", r.RTT.Milliseconds())
			}
			rows = append(rows, []string{r.Host, status, rtt})
		}

		output.Table([]string{"Host", "Status", "RTT"}, rows)
		fmt.Printf("\n%s %d online / %d total\n",
			output.Bold("Resumo:"), upCount, len(results))

		return nil
	},
}

func isCIDR(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			return true
		}
	}
	return false
}

func init() {
	pingCmd.Flags().IntVarP(&pingTimeout, "timeout", "t", 1000, "timeout por host em ms")
	pingCmd.Flags().IntVarP(&pingWorkers, "workers", "w", 50, "número de pings concorrentes")
	pingCmd.Flags().BoolVar(&pingOnlyUp, "only-up", false, "exibir apenas hosts online")
	rootCmd.AddCommand(pingCmd)
}
