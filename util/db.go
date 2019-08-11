package util

import (
	"database/sql"
	"fmt"
	"strings"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// OpenDB - open database from DATABASE_URL
// NOTICE: the url must be prefixed with mysql://
func OpenDB(source string) (*sql.DB, error) {
	// Check if URL is leading with mysql://
	// For now we only support mysql
	if !strings.HasPrefix(source, "mysql://") {
		return nil, fmt.Errorf("source:%s is invalid since we only support `mysql://` prefix", source)
	}
	var dataSource = strings.Replace(source, "mysql://", "", 1)

	// ensure mysql could be opened
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, fmt.Errorf("open db error: %s", err.Error())
	}

	return db, nil
}
