package modules

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoResetCmd = "mongo-reset"

type MongoResetModule struct {
}

type mongoExecutionOptions struct {
	url     string
	verbose bool
	help    bool
}

type MongoDocument struct {
	database string
	document string
}

func (m MongoResetModule) Execute(args []string) {
	options := mongoExecutionOptions{}
	flagset := mongoResetCreateFlagSet(&options)
	flagset.Parse(args)

	if flagset.Parsed() {
		if options.help {
			flagset.PrintDefaults()
		} else {
			mongoReset(&options)
		}
	}
}

func mongoResetCreateFlagSet(options *mongoExecutionOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(MongoResetCmd, flag.ExitOnError)
	fs.BoolVar(&options.verbose, "v", false, "Verbose")
	fs.BoolVar(&options.help, "help", false, "Help")
	fs.StringVar(&options.url, "url", "", "Mongo uri (optional. Default mongodb://localhost:27017)")
	return fs
}

func mongoReset(cmdOptions *mongoExecutionOptions) {
	mongoUri, definedUri := os.LookupEnv("APP_MONGO_URI")

	//TODO read argument if defined
	if !definedUri {
		mongoUri = "mongodb://localhost:27017"
		log.Printf("Using default mongo URI %s", mongoUri)
	}

	if cmdOptions.verbose {
		log.Printf("Cleaning documents")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("%s: %s", "Error opening connection", err)
		os.Exit(1)
	}

	remove("cnp-actions", "actions", client)
	remove("cnp-actions", "archivedActions", client)
	remove("cnp-claims", "claims", client)
	remove("cnp-coverages", "coverages", client)
	remove("cnp-customers", "persons", client)
	remove("cnp-customers", "archivedPersons", client)
	remove("cnp-documents", "documentCollections", client)
	remove("cnp-documents", "signatureCallbacks", client)
	remove("cnp-orders", "orders", client)
	remove("cnp-pedra", "amlVerifications", client)
	remove("cnp-sid", "sidHistory", client)
	remove("ppi-policies", "policies", client)
	remove("ppi-policies", "archivedPolicies", client)
	remove("pp-policies", "policies", client)
	remove("pp-policies", "archivedPolicies", client)
	remove("pp-policies", "policySyncExecutions", client)

	if cmdOptions.verbose {
		log.Printf("Reset complete")
	}
}

func remove(database string, document string, client *mongo.Client) {
	log.Printf("Removing documents from %s.%s", database, document)
	client.Database(database).Collection(document).DeleteMany(context.Background(), bson.D{})
}
