package util

import (
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
	db, err := OpenDB(source)
	if err != nil {
		return err
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
