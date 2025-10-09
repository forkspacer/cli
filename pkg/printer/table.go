package printer

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// Table provides a styled table printer
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new table with headers
func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    [][]string{},
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

// Render displays the table
func (t *Table) Render() {
	if len(t.rows) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)

	// Set header
	var headerData []interface{}
	for _, h := range t.headers {
		headerData = append(headerData, h)
	}
	table.Header(headerData...)

	// Add all rows
	for _, row := range t.rows {
		var rowData []interface{}
		for _, cell := range row {
			rowData = append(rowData, cell)
		}
		table.Append(rowData...)
	}

	table.Render()
}
