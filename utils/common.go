//go:generate goderive
package utils

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	if 'a' <= s[0] && s[0] <= 'z' {
		b := []byte(s)
		b[0] = b[0] - (byte('a') - 'A')
		return string(b)
	}
	return s
}

type TableWriter struct {
	table *tablewriter.Table
}

func NewTableWriter(w io.Writer) *TableWriter {
	t := tablewriter.NewWriter(w)
	t.SetBorders(tablewriter.Border{Left: false, Right: false, Top: false, Bottom: false})
	t.SetHeaderLine(false)
	t.SetColumnSeparator("\t")
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	return &TableWriter{table: t}
}

func (tw *TableWriter) Append(row []string) {
	tw.table.Append(row)
}

func (tw *TableWriter) Render() {
	tw.table.Render()
}

// derive-set
type Str = string
