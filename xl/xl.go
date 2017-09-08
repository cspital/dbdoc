package xl

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
)

var headers = []string{
	"Primary/Foreign Key", "Field Name",
	"Caption", "Data Type", "Field Size", "Default",
	"Nullable", "Unique", "Reference", "Notes",
}

type Field struct {
	Key       string
	Name      string
	Caption   string
	Type      string
	Size      string
	Default   string
	Nullable  string
	Unique    string
	Reference string
	Notes     string
}

type TableDescription struct {
	Name    string
	Columns []*Field
}

func New(db string) *Excel {
	return &Excel{
		filename: fmt.Sprintf("%s_DataDictionary.xlsx", db),
		wkbk:     xlsx.NewFile(),
	}
}

type Excel struct {
	filename string
	wkbk     *xlsx.File
}

func (e *Excel) DescribeTable(table *TableDescription) error {
	sheet, err := e.wkbk.AddSheet(table.Name)
	if err != nil {
		return errors.Wrapf(err, "could not make sheet for %s", table.Name)
	}
	e.writeHeaders(sheet)
	e.writeFields(sheet, table.Columns)
	return nil
}

func (e *Excel) Save() error {
	if err := e.wkbk.Save(e.filename); err != nil {
		return errors.Wrapf(err, "could not save workbook %s", e.filename)
	}
	return nil
}

func (e *Excel) writeHeaders(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	for _, h := range headers {
		row.AddCell().Value = h
	}
}

func (e *Excel) writeFields(sheet *xlsx.Sheet, fields []*Field) {
	for _, f := range fields {
		row := sheet.AddRow()
		cell(row, f.Key)
		cell(row, f.Name)
		cell(row, f.Caption)
		cell(row, f.Type)
		cell(row, f.Size)
		cell(row, f.Default)
		cell(row, f.Nullable)
		cell(row, f.Unique)
		cell(row, f.Reference)
		cell(row, f.Notes)
	}
}

func cell(row *xlsx.Row, n interface{}) {
	row.AddCell().SetValue(n)
}
