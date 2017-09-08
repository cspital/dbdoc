package doc

import (
	"fmt"
	"strings"

	"github.com/biffjutsu/dbdoc/db"
)

const (
	ForeignKey = "FOREIGN KEY"
	PrimaryKey = "PRIMARY KEY"
	Unique     = "UNIQUE"
)

func NewCache(s db.SchemaService) (*SchemaCache, error) {
	var tables []db.Table
	var columns []*db.Column
	var constraints []*db.Constraint
	var err error

	if tables, err = s.Tables(); err != nil {
		return nil, err
	}

	if columns, err = s.Columns(); err != nil {
		return nil, err
	}

	if constraints, err = s.Constraints(); err != nil {
		return nil, err
	}

	s.Close()

	return &SchemaCache{
		Tables:      tables,
		Columns:     columns,
		Constraints: constraints,
	}, nil
}

type SchemaCache struct {
	Tables      []db.Table
	Columns     []*db.Column
	Constraints []*db.Constraint
}

func (s *SchemaCache) ColumnReference(table, column string) string {
	refs := []string{}
	for _, con := range s.Constraints {
		if con.Type == ForeignKey && con.ConstrainedColumn == column && con.ConstrainedTable == table {
			refs = append(refs, fmt.Sprintf("%s.%s", con.SourceTable.String, con.SourceColumn.String))
		}
	}
	return strings.Join(refs, "; ")
}

func (s *SchemaCache) ColumnKeyTypes(table, column string) string {
	var p string
	var f string
	for _, con := range s.Constraints {
		if con.ConstrainedColumn == column && con.ConstrainedTable == table {
			switch con.Type {
			case PrimaryKey:
				p = "P"
			case ForeignKey:
				f = "F"
			default:
			}
		}
	}
	if p == "" {
		return f
	}
	if f == "" {
		return p
	}
	return "P/F"
}

func (s *SchemaCache) IsColumnUnique(table, column string) bool {
	for _, con := range s.Constraints {
		if con.Type == Unique && con.ConstrainedColumn == column && con.ConstrainedTable == table {
			return true
		}
	}
	return false
}

func (s *SchemaCache) ColumnsFor(table string) []*db.Column {
	var cols []*db.Column
	for _, col := range s.Columns {
		if col.Table == table {
			cols = append(cols, col)
		}
	}
	return cols
}
