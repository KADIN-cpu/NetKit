package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
)

// Table renderiza uma tabela simples no stdout com cabeçalho e linhas.
func Table(header []string, rows [][]string) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(header)
	t.SetBorder(false)
	t.SetColumnSeparator("")
	t.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	// Ajusta cor do cabeçalho para o número real de colunas.
	colors := make([]tablewriter.Colors, len(header))
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}
	}
	t.SetHeaderColor(colors...)
	t.AppendBulk(rows)
	t.Render()
}

// Section imprime um título de seção destacado.
func Section(title string) {
	fmt.Printf("\n%s\n", Bold(Cyan("» "+title)))
}

// Up retorna o texto de status "online" colorido.
func Up(s string) string { return Green(s) }

// Down retorna o texto de status "offline" colorido.
func Down(s string) string { return Red(s) }
