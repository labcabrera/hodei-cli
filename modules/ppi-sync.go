package modules

import (
	"flag"
	"fmt"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const PpiSyncCmd = "ppi-sync"

type PpiSyncModule struct {
}

type ppiSyncOptions struct {
	entity      string
	rsql        string
	authorities string
	help        bool
	verbose     bool
	batchSize   int
}

func (m PpiSyncModule) Execute(args []string) {
	options := ppiSyncOptions{}
	flagset := ppiSyncOptionsCreateFlagSet(&options)
	flagset.Parse(args)
	if flagset.Parsed() {
		if options.help {
			flagset.PrintDefaults()
		} else {
			ppiSync(&options)
		}
	}
}

func ppiSync(options *ppiSyncOptions) (res string, err error) {
	if options.verbose {
		log.Printf("Sending synchronization request %s: %s (%s)", options.entity, options.rsql, options.authorities)
	}
	headers := amqp.Table{
		"App-Source": "hodei-cli",
	}
	body := `{}`
	res, err = client.SendAndReceive("ppi.referential", "process", body, headers, options.verbose)
	fmt.Println(res)
	return
}

func ppiSyncOptionsCreateFlagSet(options *ppiSyncOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(CheckIbanCmd, flag.ExitOnError)
	fs.StringVar(&options.entity, "entity", "", "Entity type")
	fs.StringVar(&options.rsql, "rsql", "", "RSQL search expression")
	fs.StringVar(&options.rsql, "authorities", "", "Authorities list (coma separated)")
	fs.BoolVar(&options.verbose, "v", false, "Verbose")
	fs.BoolVar(&options.help, "help", false, "Help")
	return fs
}
