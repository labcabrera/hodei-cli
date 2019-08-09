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

type MongoResetOptions struct {
	Verbose bool
	Help    bool
}

func MongoReset(cmdOptions *MongoResetOptions) {
	//TODO read confirmation

	mongoUri, definedUri := os.LookupEnv("APP_MONGO_URI")
	if !definedUri {
		mongoUri = "mongodb://localhost:27017"
		log.Printf("Using default mongo URI %s", mongoUri)
	}

	if cmdOptions.Verbose {
		log.Printf("Cleaning documents")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("%s: %s", "Error opening connection", err)
		os.Exit(1)
	}

	collectionMap := map[string]string{
		"actions":             "cnp-actions",
		"scheduledActions":    "cnp-actions",
		"batchJobExecutions":  "cnp-actions",
		"batchJobInstances":   "cnp-actions",
		"batchSequences":      "cnp-actions",
		"batchStepExecutions": "cnp-actions",
		"persons":             "cnp-customers",
		"policies":            "ppi-policies",
		"documentCollections": "cnp-documents",
		"policyOrders":        "cnp-orders",
		"claims":              "cnp-claims",
		"policyCoverages":     "cnp-coverages",
	}
	for table, database := range collectionMap {
		log.Printf("Removing documents from %s.%s", database, table)
		client.Database(database).Collection(table).DeleteMany(context.Background(), bson.D{})
	}

	if cmdOptions.Verbose {
		log.Printf("Reset complete")
	}
}

func MongoResetFlagSet(options *MongoResetOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(MongoResetCmd, flag.ExitOnError)
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
