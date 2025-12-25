package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	//	ID      string `bson:"_id"`
	Name    string `bson:"name"`
	Age     int    `bson:"age"`
	Country string `bson:"country"`
}

func insertUsersFromFile(client *mongo.Client, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return fmt.Errorf("Error unmarshaling data: %w", err)
	}

	// Преобразуем slice структур в slice interface{}
	usersInterface := make([]interface{}, len(users))
	for i, u := range users {
		usersInterface[i] = u
	}

	collection := client.Database("test").Collection("users")
	insertResult, err := collection.InsertMany(context.TODO(), usersInterface)
	if err != nil {
		return fmt.Errorf("Insert failed: %w", err)
	}

	fmt.Printf("Inserted %d documents!\n", len(insertResult.InsertedIDs))
	return nil
}

func aggregateData(client *mongo.Client) ([]bson.M, error) {
	collection := client.Database("test").Collection("users")
	pipeline := []bson.D{
		{{"$group", bson.D{{"_id", "$country"}, {"count", bson.D{{"$sum", 1}}}}}},
	}
	cur, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("Aggregation failed: %w", err)
	}
	defer cur.Close(context.TODO())

	var results []bson.M
	for cur.Next(context.TODO()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("Decoding failed: %w", err)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("Cursor error: %w", err)
	}

	return results, nil
}

func printResults(results []bson.M) {
	for _, r := range results {
		country := r["_id"].(string)
		count := r["count"].(int32)
		fmt.Printf("%s has %v users\n", country, count)
	}
}

func main() {
	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = client.Disconnect(context.TODO()) }()

	// Загрузка данных из файла
	filename := "data.json"
	err = insertUsersFromFile(client, filename)
	if err != nil {
		log.Fatalf("Failed to load data from file: %v", err)
	}

	// Выполнение агрегирования
	results, err := aggregateData(client)
	if err != nil {
		log.Fatalf("Aggregation failed: %v", err)
	}

	printResults(results)
}
