package main

import (
	"flag"
	"log"
	"os"

	"github.com/DemoHn/al-order/app"
	"github.com/DemoHn/al-order/util"
)

var envFile string
var sqlFile string

var host string
var port int

const dbEnvKey = "DATABASE_URL"

func init() {
	flag.StringVar(&envFile, "envFile", ".env", "environment file")
	flag.StringVar(&sqlFile, "sqlFile", "", "assign migration file. only used on [migration]")
	// set default host & port
	flag.StringVar(&host, "host", "127.0.0.1", "assign host")
	flag.IntVar(&port, "port", 8080, "listen port")
}

func main() {
	flag.Parse()
	//var execArg = flag.Arg(0)
	var err error

	if envFile != "" {
		util.RegisterEnvFromFile(envFile)
	}

	//if execArg == "migration" {
	err = util.ExecMigration(os.Getenv(dbEnvKey), sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Database Migration: %s succeed", sqlFile)
	//		return
	//}

	// start server
	StartServer(host, port)
}

// StartServer - start al-order server to handle requests
func StartServer(host string, port int) {
	log.Fatal(app.New().Start(host, port))
}
