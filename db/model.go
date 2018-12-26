package db

import (
	"github.com/guregu/null"
)

// Table ...
type Table = string

// Column ...
type Column struct {
	Table      string      `db:"TABLE_NAME"`
	Name       string      `db:"COLUMN_NAME"`
	IsNullable string      `db:"IS_NULLABLE"`
	Type       string      `db:"DATA_TYPE"`
	MaxLength  null.Int    `db:"CHARACTER_MAXIMUM_LENGTH"`
	Default    null.String `db:"COLUMN_DEFAULT"`
	IsIdentity bool        `db:"IS_IDENTITY"`
}

// Constraint ...
type Constraint struct {
	Name              string      `db:"CONSTRAINT_NAME"`
	Type              string      `db:"CONSTRAINT_TYPE"`
	ConstrainedTable  string      `db:"CONSTRAINED_TABLE"`
	ConstrainedColumn string      `db:"CONSTRAINED_COLUMN"`
	SourceTable       null.String `db:"SOURCE_TABLE"`
	SourceColumn      null.String `db:"SOURCE_COLUMN"`
}
