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

const usage = `usage: dbdoc.exe -server=mydb.server.org -db=myDatabase

args:
	-server       database server to connect to
	-db           database to make scaffolding for
	-help         see this
`

func main() {
	server := flag.String("server", "", "server to connect to")
	database := flag.String("db", "", "database to document")
	help := flag.Bool("help", false, "help")
	flag.Parse()

	if *help {
		fmt.Println(usage)
		os.Exit(0)
	}

	opts, err := config.Validate(*server, *database)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	schemaService, err := db.New(opts)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cache, err := doc.NewCache(schemaService)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	describer := doc.NewDescriber(cache, xl.New(opts.Database))
	err = describer.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
