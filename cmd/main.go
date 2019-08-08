package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var envFile string
var sqlFile string

const dbEnvKey = "DATABASE_URL"

func init() {
	flag.StringVar(&envFile, "envFile", ".env", "environment file")
	flag.StringVar(&sqlFile, "sqlFile", "", "assign migration file. only used on [migration]")
}

func main() {
	flag.Parse()
	var execArg = flag.Arg(0)
	var err error

	if envFile != "" {
		RegisterEnvFromFile(envFile)
	}

	if execArg == "migration" {
		err = ExecMigration(os.Getenv(dbEnvKey), sqlFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// start server
}

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

	return nil
}

// RegisterEnvFromFile - register env data from file
func RegisterEnvFromFile(envFile string) error {
	envs, err := ioutil.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("read env file error: read file(%s) failed", envFile)
	}
	re, _ := regexp.Compile("^([A-Z_]+)=(.*)$")

	var strEnvs = strings.Split(string(envs), "\n")
	for _, strEnv := range strEnvs {
		keyValuePair := re.FindStringSubmatch(strEnv)
		os.Setenv(keyValuePair[1], keyValuePair[2])
	}

	return nil
}

// StartServer - start al-order server to handle requests
func StartServer(host string, port int) {

}
