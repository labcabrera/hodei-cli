package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/labcabrera/hodei-cli/modules"
)

const version = "0.7.0-SNAPSHOT"
const versionCmd = "version"

func main() {

	var moduleMap map[string]modules.HodeiCliModule
	moduleMap = make(map[string]modules.HodeiCliModule)
	moduleMap[modules.ListScheduledActionsCmd] = modules.ListScheduledActionsModule{}
	moduleMap[modules.MongoResetCmd] = modules.MongoResetModule{}
	moduleMap[modules.CheckIbanCmd] = modules.CheckIbanModule{}
	moduleMap[modules.PpiSyncCmd] = modules.PpiSyncModule{}
	moduleMap[modules.CustomerSearchCmd] = modules.CustomerSearchModule{}
	moduleMap[modules.SignatureRequestCmd] = modules.SignatureRequestModule{}

	if len(os.Args) < 2 {
		usage(moduleMap)
		return
	}

	cmd := os.Args[1]
	module, check := moduleMap[cmd]
	if !check {
		usage(moduleMap)
		os.Exit(1)
	} else {
		module.Execute(os.Args[2:])
		os.Exit(0)
	}
}

func usage(moduleMap map[string]modules.HodeiCliModule) {
	fmt.Println(`
Usage: hodei-cli COMMAND [OPTIONS]")

Commands:`)

	commands := make([]string, len(moduleMap))
	index := 0
	for k, _ := range moduleMap {
		commands[index] = k
		index = index + 1
	}
	sort.Strings(commands)
	for _, cmd := range commands {
		fmt.Printf("  %s\n", cmd)
	}

}
