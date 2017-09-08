package db

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/biffjutsu/dbdoc/config"
	// I be the driver
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

const dbstring = "server=%s;database=%s;Trusted_Connection=True;Application Name=db-doc"

// SchemaService ...
// Encapsulates data retrieval pieces.
type SchemaService interface {
	Tables() ([]Table, error)
	Columns() ([]*Column, error)
	Constraints() ([]*Constraint, error)
	Close() error
}

// New ...
// Creates a new schema data service.
func New(o config.Options) (SchemaService, error) {
	db, err := sqlx.Connect("mssql", fmt.Sprintf(dbstring, o.Server, o.Database))
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to %s.%s", o.Server, o.Database)
	}
	return &sqlService{db}, nil
}

type sqlService struct {
	db *sqlx.DB
}

func (s *sqlService) Tables() ([]Table, error) {
	tables := []Table{}

	err := s.db.Select(&tables, tableQuery)
	if err != nil {
		return nil, errors.Wrap(err, "could not query table list")
	}
	return tables, nil
}

func (s *sqlService) Columns() ([]*Column, error) {
	columns := []*Column{}

	err := s.db.Select(&columns, columnQuery)
	if err != nil {
		return nil, errors.Wrap(err, "could not query column list")
	}
	return columns, nil
}

func (s *sqlService) Constraints() ([]*Constraint, error) {
	ctrs := []*Constraint{}

	err := s.db.Select(&ctrs, constraintQuery)
	if err != nil {
		return nil, errors.Wrap(err, "could not query contraint list")
	}
	return ctrs, nil
}

func (s *sqlService) Close() error {
	return s.db.Close()
}

const tableQuery = `SELECT TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_NAME != 'sysdiagrams' AND
TABLE_TYPE = 'BASE TABLE';`

const columnQuery = `SELECT c.TABLE_NAME, c.COLUMN_NAME, c.IS_NULLABLE,
c.DATA_TYPE, c.CHARACTER_MAXIMUM_LENGTH, c.COLUMN_DEFAULT,
CAST(COLUMNPROPERTY(object_id(c.TABLE_NAME), c.COLUMN_NAME, 'IsIdentity') as bit) IS_IDENTITY
FROM INFORMATION_SCHEMA.COLUMNS c
INNER JOIN INFORMATION_SCHEMA.TABLES t on c.TABLE_NAME = t.TABLE_NAME
WHERE t.TABLE_NAME != 'sysdiagrams' AND TABLE_TYPE = 'BASE TABLE'
ORDER BY c.TABLE_NAME, c.ORDINAL_POSITION;`

const constraintQuery = `SELECT t.CONSTRAINT_NAME, t.CONSTRAINT_TYPE, t.TABLE_NAME CONSTRAINED_TABLE,
k.COLUMN_NAME CONSTRAINED_COLUMN, f.TABLE_NAME SOURCE_TABLE, f.COLUMN_NAME SOURCE_COLUMN
FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS t
INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE k on t.CONSTRAINT_NAME = k.CONSTRAINT_NAME
LEFT JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS r on t.CONSTRAINT_NAME = r.CONSTRAINT_NAME
LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE f on r.UNIQUE_CONSTRAINT_NAME = f.CONSTRAINT_NAME
WHERE t.CONSTRAINT_TYPE IN ('PRIMARY KEY', 'FOREIGN KEY', 'UNIQUE') AND
t.TABLE_NAME != 'sysdiagrams';`
