package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/KADIN-cpu/netkit/internal/output"
	"github.com/KADIN-cpu/netkit/internal/portscan"
)

var (
	psTimeout int
	psWorkers int
	psPorts   string
	psAll     bool
)

var portscanCmd = &cobra.Command{
	Use:     "portscan <host>",
	Aliases: []string{"ports", "scan"},
	Short:   "Verifica portas TCP abertas em um host",
	Long: `Verifica quais portas TCP estão abertas em um host.

Por padrão testa uma lista de portas comuns. Use --ports para
especificar portas ou faixas, ou --all para 1-65535.

Exemplos:
  netkit portscan 192.168.0.10
  netkit portscan 192.168.0.10 --ports 22,80,443
  netkit portscan 192.168.0.10 --ports 1-1024
  netkit portscan 192.168.0.10 --all`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		timeout := time.Duration(psTimeout) * time.Millisecond

		var ports []int
		switch {
		case psAll:
			ports = portscan.PortRange(1, 65535)
		case psPorts != "":
			p, err := parsePorts(psPorts)
			if err != nil {
				return err
			}
			ports = p
		default:
			ports = portscan.CommonPorts()
		}

		output.Section(fmt.Sprintf("Scan de %d porta(s) em %s", len(ports), host))
		results := portscan.Scan(host, ports, timeout, psWorkers)

		rows := make([][]string, 0)
		openCount := 0
		for _, r := range results {
			if !r.Open {
				continue
			}
			openCount++
			rows = append(rows, []string{
				fmt.Sprintf("%d", r.Port),
				output.Up("aberta"),
				r.Service,
			})
		}

		if openCount == 0 {
			fmt.Println(output.Yellow("Nenhuma porta aberta encontrada."))
			return nil
		}

		output.Table([]string{"Porta", "Status", "Serviço"}, rows)
		fmt.Printf("\n%s %d porta(s) aberta(s) de %d testada(s)\n",
			output.Bold("Resumo:"), openCount, len(ports))
		return nil
	},
}

// parsePorts interpreta uma string tipo "22,80,443" ou "1-1024" ou misto.
func parsePorts(s string) ([]int, error) {
	seen := make(map[int]bool)
	var ports []int

	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return nil, fmt.Errorf("porta inválida: %q", bounds[0])
			}
			end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err != nil {
				return nil, fmt.Errorf("porta inválida: %q", bounds[1])
			}
			for _, p := range portscan.PortRange(start, end) {
				if !seen[p] {
					seen[p] = true
					ports = append(ports, p)
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("porta inválida: %q", part)
			}
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("porta fora do intervalo: %d", p)
			}
			if !seen[p] {
				seen[p] = true
				ports = append(ports, p)
			}
		}
	}
	return ports, nil
}

func init() {
	portscanCmd.Flags().IntVarP(&psTimeout, "timeout", "t", 500, "timeout por porta em ms")
	portscanCmd.Flags().IntVarP(&psWorkers, "workers", "w", 100, "número de conexões concorrentes")
	portscanCmd.Flags().StringVarP(&psPorts, "ports", "p", "", "portas a testar (ex: 22,80,443 ou 1-1024)")
	portscanCmd.Flags().BoolVar(&psAll, "all", false, "testar todas as portas (1-65535)")
	rootCmd.AddCommand(portscanCmd)
}
