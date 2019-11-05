package main

import (
	"fmt"
	"os"

	"github.com/labcabrera/hodei-cli/modules"
)

const version = "0.7.0-SNAPSHOT"
const versionCmd = "version"

func main() {

	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := os.Args[1]

	var moduleMap map[string]modules.HodeiCliModule
	moduleMap = make(map[string]modules.HodeiCliModule)

	moduleMap[modules.ListScheduledActionsCmd] = modules.ListScheduledActionsModule{}
	moduleMap[modules.MongoResetCmd] = modules.MongoResetModule{}
	moduleMap[modules.CheckIbanCmd] = modules.CheckIbanModule{}
	moduleMap[modules.PpiSyncCmd] = modules.PpiSyncModule{}
	moduleMap[modules.CustomerSearchCmd] = modules.CustomerSearchModule{}
	moduleMap[modules.SignatureRequestCmd] = modules.SignatureRequestModule{}

	module, check := moduleMap[cmd]

	if !check {
		usage()
		os.Exit(1)
	} else {
		module.Execute(os.Args[2:])
		os.Exit(0)
	}
}

func usage() {
	fmt.Println(`
Usage: hodei-cli COMMAND [OPTIONS]")

Commands:
  ` + modules.CustomerSearchCmd + `
  ` + modules.CheckIbanCmd + `
  ` + modules.MongoResetCmd + `
  ` + modules.SignatureRequestCmd + `
  ` + modules.PpiSyncCmd + `
  ` + versionCmd)
}
