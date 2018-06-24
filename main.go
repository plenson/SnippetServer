// Copyright Peter Lenson.
// All Rights Reserved

// This package implements an example snippet service.
package main

import (
	dbpkg "github.com/textioHQ/interview-peter-lenson/database"
	pspkg "github.com/textioHQ/interview-peter-lenson/processargs"
	rtpkg "github.com/textioHQ/interview-peter-lenson/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get commandline arguments
	argParms := pspkg.ProcessArgs()

	// Setup snippet database
	//#IF BUILD BOLT
	sniptDB, err := dbpkg.SetupDatabaseBolt(argParms)
	//#IF BUILD BOW
	//sniptDB, err := dbpkg.SetupDatabaseBow(argParms)
	if err != nil {
		log.Fatalf("Could not create database! %s", err)
		os.Exit(1)
	}

	// Initialize route connection to snippet database
	//#IF BUILD BOLT
	rtpkg.InitBolt(*sniptDB)
	//#IF BUILD BOW
	//rtpkg.InitBow(*sniptDB)

	// Listen for requests and dispatch to appropriate route for handling
	//  gzip compress all responses,logger is called first on each route request
	log.Fatal(http.ListenAndServe(":"+argParms.Port, rtpkg.SetUpRouteHandlers()))
	//		log.Fatal(http.ListenAndServe(":"+argParms.Port, handlers.LoggingHandler(os.Stdout,rtpkg.SetUpRouteHandlers())))
	//	handlers.CompressHandler(handlers.LoggingHandler(os.Stdout, rtpkg.SetUpRouteHandlers()))))
}
