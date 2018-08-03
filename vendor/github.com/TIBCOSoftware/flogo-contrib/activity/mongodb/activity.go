package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-mongodb")

const (
	methodGet     = "GET"
	methodDelete  = "DELETE"
	methodInsert  = "INSERT"
	methodReplace = "REPLACE"
	methodUpdate  = "UPDATE"

	ivConnectionURI = "uri"
	ivDbName        = "dbName"
	ivCollection    = "collection"
	ivMethod        = "method"

	ivKeyName  = "keyName"
	ivKeyValue = "keyValue"
	ivData     = "data"

	ovOutput = "output"
	ovCount = "count"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

/*
Integration with MongoDb
inputs: {uri, dbName, collection, method, [keyName, keyValue, value]}
outputs: {output, count}
*/
type MongoDbActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MongoDbActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *MongoDbActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - MongoDb integration
func (a *MongoDbActivity) Eval(ctx activity.Context) (done bool, err error) {

	//mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
	connectionURI, _ := ctx.GetInput(ivConnectionURI).(string)
	dbName, _ := ctx.GetInput(ivDbName).(string)
	collectionName, _ := ctx.GetInput(ivCollection).(string)
	method, _ := ctx.GetInput(ivMethod).(string)
	keyName, _ := ctx.GetInput(ivKeyName).(string)
	keyValue, _ := ctx.GetInput(ivKeyValue).(string)
	value := ctx.GetInput(ivData)

	//todo implement shared sessions
	client, err := mongo.NewClient(connectionURI)
	if err != nil {
		activityLog.Errorf("Connection error: %v", err)
		return false, err
	}

	db := client.Database(dbName)

	coll := db.Collection(collectionName)

	switch strings.ToUpper(method) {
	case methodGet:
		result := coll.FindOne(context.Background(), bson.NewDocument(bson.EC.String(keyName, keyValue)))
		val := make(map[string]interface{})
		err := result.Decode(val)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Get Results $#v", result)

		ctx.SetOutput(ovOutput, val)
	case methodDelete:
		result, err := coll.DeleteMany(
			context.Background(),
			bson.NewDocument(
				bson.EC.String(keyName, keyValue),
			),
		)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Delete Results $#v", result)

		ctx.SetOutput(ovCount, result.DeletedCount)
	case methodInsert:
		result, err := coll.InsertOne(
			context.Background(),
			value,
		)
		if err != nil {
			return false, err
		}
		activityLog.Debugf("Insert Results $#v", result)

		ctx.SetOutput(ovOutput, result.InsertedID)
	case methodReplace:
		result, err := coll.ReplaceOne(
			context.Background(),
			bson.NewDocument(
				bson.EC.String(keyName, keyValue),
			),
			value,
		)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Replace Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		ctx.SetOutput(ovCount, result.ModifiedCount)

	case methodUpdate:
		result, err := coll.UpdateOne(
			context.Background(),
			bson.NewDocument(
				bson.EC.String(keyName, keyValue),
			),
			bson.NewDocument(
				bson.EC.Interface("$set", value),
			),
		)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Update Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		ctx.SetOutput(ovCount, result.ModifiedCount)
	default:
		activityLog.Errorf("unsupported method '%s'", method)
		return false, fmt.Errorf("unsupported method '%s'", method)
	}

	return true, nil
}
