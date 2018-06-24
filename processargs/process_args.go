// Copyright Peter Lenson.
// All Rights Reserved

package processargs

import (
	"flag"
	"fmt"
	cmnpkg "github.com/plenson/SnippetService/common"
	"log"
	"os"
	"path"
)

const version = "v0.0.0"

func usage(appName, version string) {
	fmt.Printf("\nOPTIONS:\n\n")
	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > 1 {
			fmt.Printf("    -%s, -%s\t%s\n", f.Name[0:1], f.Name, f.Usage)
		}
	})
	fmt.Printf("\n\nVersion %s\n", version)
}

// Processes command line arguments
func ProcessArgs() cmnpkg.DbParams {

	dbP := new(cmnpkg.DbParams)
	flag.BoolVar(&dbP.ShowHelp, "h", false, "display help")
	flag.BoolVar(&dbP.ShowHelp, "help", false, "display help")

	appName := path.Base(os.Args[0])
	flag.Parse()

	args := flag.Args()

	if dbP.ShowHelp == true {
		usage(appName, version)
		os.Exit(0)
	}
	// If non-flag options are included assume bolt db is specified.
	if len(args) > 0 {
		dbP.DbName = args[0]
	}

	if dbP.DbName == "" {
		usage(appName, version)
		log.Printf("\nERROR: Missing boltdb name\n")
		os.Exit(1)
	}
	if len(args) > 1 {
		dbP.HmName = args[1]
	}

	if dbP.HmName == "" {
		usage(appName, version)
		log.Printf("\nERROR: Missing hashmap name\n")
		os.Exit(1)
	}
	if len(args) > 2 {
		dbP.Port = args[2]
	}

	if dbP.Port == "" {
		usage(appName, version)
		log.Printf("\nERROR: Missing port name\n")
		os.Exit(1)
	}

	if len(args) > 3 {
		dbP.DataVolPath = args[3]
	}
	if dbP.DataVolPath == "" {
		usage(appName, version)
		log.Printf("\nERROR: Missing data volume path \n")
		os.Exit(1)
	}
	return *dbP
}
