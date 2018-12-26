package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/biffjutsu/dbdoc/config"
	"github.com/biffjutsu/dbdoc/db"
	"github.com/biffjutsu/dbdoc/doc"
	"github.com/biffjutsu/dbdoc/xl"
)

const usage = `
usage: dbdoc.exe -server=mydb.server.org -db=myDatabase

args:
	-server       required. database server to connect to
	-db           required. database to make scaffolding for
	-uid          optional. user to connect with, if omitted will attempt to use trusted connection
	-pw           optional. passord to connect with
	-help         see this
`

func unwrap(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(opts config.Options) {
	schemaService, err := db.New(opts)
	unwrap(err)

	cache, err := doc.NewCache(schemaService)
	unwrap(err)

	describer := doc.NewDescriber(cache, xl.New(opts.Database))
	err = describer.Run()
	unwrap(err)
}

func main() {
	server := flag.String("server", "", "server to connect to")
	database := flag.String("db", "", "database to document")
	username := flag.String("uid", "", "user to connect as")
	password := flag.String("pw", "", "password")
	help := flag.Bool("help", false, "help")
	flag.Parse()

	if *help || (*server == "" || *database == "") {
		fmt.Println(usage)
		os.Exit(0)
	}

	opts, err := config.Validate(*server, *database, *username, *password)
	unwrap(err)

	run(opts)
	os.Exit(0)
}
