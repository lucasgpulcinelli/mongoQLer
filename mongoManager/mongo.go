package mongoManager

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Login logs in to a mongoDB database, returning the connection.
func Login(url, databaseName, user, password string) (*mongo.Database, error) {
	credential := options.Credential{
		AuthSource: databaseName,
		Username:   user,
		Password:   password,
	}
	clientOpts := options.Client().ApplyURI("mongodb://" + url).
		SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}

	db := client.Database(databaseName)

	var result bson.M
	err = db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return db, nil
}
