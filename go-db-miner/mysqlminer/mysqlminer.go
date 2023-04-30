package mysqlminer

import (
	"database/sql"
	"fmt"
	dbminer "go-db-miner/dbminer"
)

type MySQLMiner struct {
	Host string
	Db   sql.DB
}

func New(host string) (*MySQLMiner, error) {
	m := MySQLMiner{Host: host}

	err := m.connect()
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *MySQLMiner) connect() error {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("root:password@tcp(%s:3306)/information_schema", m.Host),
	)
	if err != nil {
		return err
	}

	m.Db = *db
	return nil
}

func (m *MySQLMiner) GetSchema() (*dbminer.Schema, error) {

	var s = new(dbminer.Schema)
	sql := `SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM columns 
			WHERE TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
			ORDER BY TABLE_SCHEMA, TABLE_NAME`

	schema_rows, err := m.Db.Query(sql)

	if err != nil {
		return nil, err
	}

	defer schema_rows.Close()

	var prev_schema, prev_table string

	var db dbminer.Database
	var table dbminer.Table

	for schema_rows.Next() {
		var curr_schema, curr_table, curr_col string
		if err := schema_rows.Scan(&curr_schema, &curr_table, &curr_col); err != nil {
			return nil, err
		}

		if curr_schema != prev_schema {
			if prev_schema != "" {
				db.Tables = append(db.Tables, table)
				s.Databases = append(s.Databases, db)
			}
			db = dbminer.Database{Name: curr_schema, Tables: []dbminer.Table{}}
			prev_schema = curr_schema
			prev_table = ""
		}

		if curr_table != prev_table {
			if prev_table != "" {
				db.Tables = append(db.Tables, table)
			}
			table = dbminer.Table{Name: curr_table, Columns: []string{}}
			prev_table = curr_table
		}
		table.Columns = append(table.Columns, curr_col)

	}

	db.Tables = append(db.Tables, table)
	s.Databases = append(s.Databases, db)
	if err := schema_rows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}
