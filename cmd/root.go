package cmd

import (
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "netkit",
	Short: "netkit - toolkit de infraestrutura e rede para analistas de TI",
	Long: `netkit é uma ferramenta de linha de comando multiplataforma para
diagnóstico de sistemas e rede.

Reúne em um único binário as tarefas mais comuns do dia a dia de TI:
informações do sistema, varredura de rede local, ping em massa e
checagem de portas — sem precisar de scripts soltos.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute roda o comando raiz.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate("netkit v{{.Version}}\n")
}
