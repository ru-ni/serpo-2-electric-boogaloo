package main

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
)

func makeTable(headers, footers []string, payload [][]string) string {
	if len(headers) != len(footers) || len(payload[0]) != len(headers) {
		//data mismatch, return error string
		return "(ノ`´)ノ ~┻━┻ "
	} //should have full data alignment from here on

	var b bytes.Buffer
	data := [][]string{}

	table := tablewriter.NewWriter(&b)
	table.SetHeader(headers)

	for _, rowdata := range payload {
		data = append(data, rowdata)
	}

	table.SetFooter(footers) // Add Footer
	table.SetBorder(false)

	table.AppendBulk(data) // Add Bulk Data
	table.SetAlignment(1)
	table.Render()
	return fmt.Sprintf("```%v```", &b)
}
