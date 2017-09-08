package doc

import (
	"strconv"

	"github.com/biffjutsu/dbdoc/db"
	"github.com/biffjutsu/dbdoc/xl"
)

func NewDescriber(s *SchemaCache, e *xl.Excel) *DBDoc {
	return &DBDoc{
		cache: s,
		excel: e,
	}
}

type DBDoc struct {
	cache *SchemaCache
	excel *xl.Excel
}

func (d *DBDoc) Run() error {
	for _, t := range d.cache.Tables {
		desc := xl.TableDescription{
			Name:    t,
			Columns: make([]*xl.Field, 0),
		}

		for _, col := range d.cache.ColumnsFor(t) {
			field := d.makeField(t, col)
			desc.Columns = append(desc.Columns, field)
		}

		err := d.excel.DescribeTable(&desc)
		if err != nil {
			return err
		}
	}
	return d.excel.Save()
}

func (d *DBDoc) makeField(table string, col *db.Column) *xl.Field {
	field := new(xl.Field)
	field.Name = col.Name
	field.Type = col.Type
	field.Nullable = col.IsNullable

	if col.MaxLength.Valid {
		field.Size = strconv.FormatInt(col.MaxLength.Int64, 10)
	} else {
		field.Size = ""
	}

	if col.Default.Valid {
		field.Default = col.Default.String
	} else {
		field.Default = ""
	}

	if d.cache.IsColumnUnique(table, col.Name) {
		field.Unique = "YES"
	} else {
		field.Unique = "NO"
	}

	field.Key = d.cache.ColumnKeyTypes(table, col.Name)
	field.Reference = d.cache.ColumnReference(table, col.Name)
	if col.IsIdentity {
		field.Caption = "Identity"
	}
	return field
}
