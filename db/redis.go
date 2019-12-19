package db

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/go-redis/redis/v7"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // Client is a wrapping struct for the redis client to add functionality
// type Client struct {
// 	RedisClient     *redis.Client
// 	MongoClient     *mongo.Client
// 	MongoCollection *mongo.Collection
// }

// // GetClient gets a new database client
// func GetClient(address, password string, db int) (*Client, error) {
// 	client := Client{}
// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr:     address,
// 		Password: password,
// 		DB:       db,
// 	})
// 	client.RedisClient = redisClient

// 	_, err := client.RedisClient.Ping().Result()
// 	if err != nil {
// 		fmt.Println("Got error on connecting to redis " + err.Error())
// 		return &client, err
// 	}

// 	return &client, nil
// }

// // SetValueString sets a string value for the provided key in the database
// func (client *Client) SetValueString(key, value string) error {
// 	return client.RedisClient.Set(key, value, 0).Err()
// }

// // SetValueBytes converts the provided byte array into a string and saves it to the database
// func (client *Client) SetValueBytes(key string, value []byte) error {
// 	return client.SetValueString(key, string(value))
// }

// // SetValueStruct takes any struct, marshals it and saves it to the db as a astring
// func (client *Client) SetValueStruct(key string, value interface{}) error {
// 	byteArray, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}

// 	return client.SetValueBytes(key, byteArray)
// }

// // GetValue returns the value as string matching the provided key
// func (client *Client) GetValue(key string) (string, error) {
// 	return client.RedisClient.Get(key).Result()
// }
