package modules

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labcabrera/hodei-cli/model"
)

const ListScheduledActionsCmd = "scheduled-actions"
const printTemplate = "|%-25s|%-25s|%-10s|%-25s|%-20s|%-5s|\n"

type ListScheduledActionsModule struct {
}

type listScheduledActionsOptions struct {
	actionType string
	executed   bool
	verbose    bool
	help       bool
}

func (m ListScheduledActionsModule) Execute(args []string) {
	executionOptions := listScheduledActionsOptions{}
	flagset := listScheduledActionsCreateFlagSet(&executionOptions)
	flagset.Parse(args)

	if flagset.Parsed() {
		if executionOptions.help {
			flagset.PrintDefaults()
		} else {
			listScheduledActions(&executionOptions)
		}
	}
}

func listScheduledActionsCreateFlagSet(executionOptions *listScheduledActionsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(ListScheduledActionsCmd, flag.ExitOnError)
	fs.BoolVar(&executionOptions.verbose, "v", false, "Verbose")
	fs.BoolVar(&executionOptions.help, "help", false, "Help")
	fs.BoolVar(&executionOptions.executed, "executed", false, "Executed (optional)")
	fs.StringVar(&executionOptions.actionType, "type", "", "Action type (optional)")
	return fs
}

func listScheduledActions(executionOptions *listScheduledActionsOptions) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	if executionOptions.verbose {
		fmt.Println("Connected to MongoDB")
	}

	collection := client.Database("cnp-actions").Collection("scheduledActions")

	//TODO
	//filter := bson.D{{"entityType", "document"}}
	filter := bson.D{{}}
	//if executionOptions.executed {
	//	filter = bson.D{{"executed", "null"}}
	//} else {
	//	filter = bson.D{{"executed", "{$ne, null"}}
	//}

	findOptions := options.Find()
	findOptions.SetLimit(25)

	var results []*model.ScheduledAction
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for cur.Next(context.TODO()) {
		var elem model.ScheduledAction
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		results = append(results, &elem)
	}
	cur.Close(context.TODO())

	fmt.Printf(printTemplate, "Id", "EntityId", "EntityType", "ActionType", "Execution", "Code")
	for _, action := range results {
		executed := ""
		if !action.Executed.IsZero() {
			executed = action.Executed.Format("2006-01-02 15:04:05")
		}
		fmt.Printf(printTemplate, action.Id.Hex(), action.EntityId, action.EntityType, action.ActionType, executed, action.Result.Code)
	}

}
