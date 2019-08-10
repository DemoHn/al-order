package util

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// migration.go: all migration utils to help exec migration process

// ExecMigration - execute sql migration
// param: source [string]: available db source, must be `mysql://` URL format
// param: file [string]: SQL file to execute
func ExecMigration(source string, file string) error {
	// Check if URL is leading with mysql://
	// For now we only support mysql
	if !strings.HasPrefix(source, "mysql://") {
		return fmt.Errorf("source:%s is invalid since we only support `mysql://` prefix", source)
	}
	var dataSource = strings.Replace(source, "mysql://", "", 1)

	// ensure mysql could be opened
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return fmt.Errorf("open db error: %s", err.Error())
	}

	// ensure file could be opened
	stmt, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("open file error: %s", err.Error())
	}

	requests := strings.Split(string(stmt), ";")
	for _, request := range requests {
		if request == "" {
			continue
		}
		_, e := db.Exec(request)
		if e != nil {
			return fmt.Errorf("exec sql statement error: %s", e.Error())
		}
	}

	return db.Close()
}
