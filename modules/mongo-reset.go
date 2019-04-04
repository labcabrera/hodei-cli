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
		log.Printf("Cleaning mongo documents")
	}

	if cmdOptions.Verbose {
		log.Printf("Cleaning mongo documents (%s)", mongoUri)
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("%s: %s", "Error opening connection", err)
		os.Exit(1)
	}

	client.Database("cnp-actions").Collection("entityActions").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-commons").Collection("batchJobExecutions").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-commons").Collection("batchJobInstances").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-commons").Collection("batchSequences").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-commons").Collection("batchStepExecutions").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-commons").Collection("batchStepExecutions").DeleteMany(context.Background(), bson.D{})
	//TODO dont remove networks after referential sync changes
	client.Database("cnp-commons").Collection("networks").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-customers").Collection("legalEntities").DeleteMany(context.Background(), bson.D{})
	client.Database("cnp-customers").Collection("persons").DeleteMany(context.Background(), bson.D{})
	client.Database("ppi-policies").Collection("policies").DeleteMany(context.Background(), bson.D{})

	if cmdOptions.Verbose {
		log.Printf("Operacion complete")
	}

}

func MongoResetFlagSet(options *MongoResetOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(MongoResetCmd, flag.ExitOnError)
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
