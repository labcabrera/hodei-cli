package modules

import (
	"flag"
	"fmt"
)

const ListScheduledActionsCmd = "scheduled-actions"

type ListScheduledActionsModule struct {
}

type listScheduledActionsOptions struct {
	actionType string
	verbose    bool
	help       bool
}

func (m ListScheduledActionsModule) Execute(args []string) {
	options := listScheduledActionsOptions{}
	flagset := listScheduledActionsCreateFlagSet(&options)
	flagset.Parse(args)

	if flagset.Parsed() {
		if options.help {
			flagset.PrintDefaults()
		} else {
			listScheduledActions(&options)
		}
	}
}

func listScheduledActionsCreateFlagSet(options *listScheduledActionsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(ListScheduledActionsCmd, flag.ExitOnError)
	fs.BoolVar(&options.verbose, "v", false, "Verbose")
	fs.BoolVar(&options.help, "help", false, "Help")
	fs.StringVar(&options.actionType, "type", "", "Action type (optional)")
	return fs
}

func listScheduledActions(options *listScheduledActionsOptions) {
	//TODO not implemented
	fmt.Println("Not implemented")
}
