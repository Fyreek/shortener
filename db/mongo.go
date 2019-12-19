package db

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/fyreek/shortener/logging"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (mDB *MongoDB) Connect(ip string, port, timeout int) error {
	cAddress := "mongodb://" + ip + ":" + strconv.Itoa(port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cAddress))
	if err != nil {
		return err
	}
	mDB.Client = client
	return nil
}

func (mDB *MongoDB) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := mDB.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}
	return true
}

func (mDB *MongoDB) SetDatabase(name string) {
	mDB.Database = mDB.Client.Database(name)
}

func (mDB *MongoDB) GetSingleEntry(collection, column, value string, iStruct interface{}) error {
	coll := mDB.Database.Collection(collection)
	filter := bson.M{column: value}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := coll.FindOne(ctx, filter).Decode(iStruct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoDocument
		}
		return err
	}
	return nil
}

func (mDB *MongoDB) GetMultipleEntries(collection, column, value string, sort map[string]interface{}, limit int) ([][]byte, error) {
	iLimit := int64(limit)

	coll := mDB.Database.Collection(collection)
	opt := options.Find()
	opt.Sort = sort
	opt.Limit = &iLimit
	filter := bson.M{column: value}
	if value == "" {
		filter = bson.M{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := coll.Find(ctx, filter, opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoDocument
		}
		return nil, err
	}

	defer cur.Close(ctx)

	byteAArray := make([][]byte, 0)
	for cur.Next(ctx) {
		elem := &bson.D{}
		if err := cur.Decode(elem); err != nil {
			return byteAArray, err
		}
		byteElem, err := bson.MarshalExtJSON(elem, false, true)
		if err != nil {
			return byteAArray, err
		}
		byteAArray = append(byteAArray, byteElem)
	}
	return byteAArray, nil
}

func (mDB *MongoDB) InsertSingleEntry(collection string, value interface{}) error {
	coll := mDB.Database.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	insertResult, err := coll.InsertOne(ctx, value)
	if err != nil {
		return err
	}
	logging.Log(logging.Debug, "Inserted doc with id:", insertResult.InsertedID)
	return nil
}

func (mDB *MongoDB) UpdateSingleEntry(collection, filterColumn, filterValue string, values interface{}) error {
	coll := mDB.Database.Collection(collection)
	filter := bson.M{filterColumn: filterValue}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{"$set", values},
	}
	updateResult, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	logging.Log(logging.Debug, "Updated entry:", updateResult.MatchedCount)
	return nil
}

func (mDB *MongoDB) DeleteSingleEntry(collection, filterColumn, filterValue string) error {
	coll := mDB.Database.Collection(collection)
	filter := bson.M{filterColumn: filterValue}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	logging.Log(logging.Debug, "Deleted entry:", res.DeletedCount)
	return nil
}
