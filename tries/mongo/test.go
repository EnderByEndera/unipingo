package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func save(ctx context.Context, client *mongo.Client) (err error) {
	type Post struct {
		Title string `bson:"title,omitempty"`
		Body  string `bson:"body,omitempty"`
	}

	type AAA struct {
		AAA   string `bson:"aaa,omitempty"`
		Posts []Post `bson:"post,omitempty"`
	}

	/*
	   Get my collection instance
	*/
	collection := client.Database("blog").Collection("posts")

	/*
	   Insert documents
	*/
	docs := []interface{}{
		// bson.D
		bson.D{{"title", "World"}, {"body", "Hello World"}},
		bson.D{{"title", "Mars"}, {"body", "Hello Mars"}},
		bson.D{{"title", "Pluto"}, {"body", "Hello Pluto"}},
	}

	byts, err := bson.Marshal(AAA{"xx", []Post{{"aaa", "aaaa"}}})
	fmt.Println(byts)
	if err != nil {
		log.Fatal(err)
	}
	postsUnmarshalled := AAA{}

	err = bson.Unmarshal(byts, &postsUnmarshalled)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("unmarshalled", postsUnmarshalled)
	res, insertErr := collection.InsertMany(ctx, docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)
	/*
	   Iterate a cursor and print it
	*/
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var posts []Post
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}
	fmt.Println(posts)
	return
}

func main() {

	/*
	   Connect to my cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:10200"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	save(ctx, client)
}
