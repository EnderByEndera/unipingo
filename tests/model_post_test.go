package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestModelPosts(t *testing.T) {
	byts, _ := ioutil.ReadFile("post.json")
	post := models.Post{}
	err := json.Unmarshal(byts, &post)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	res, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(string(res))

	mongoConn := db.GetMongoConn()
	client := mongoConn.Client

	collection := client.Database("blog").Collection("posts")

	// insertedDocID, err := primitive.ObjectIDFromHex("63f7793cd11f8c5c7258a06d")
	insertedDocID, err := services.GetPostsService().NewPost(&post)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	postGot, err := services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if postGot.DocumentID != insertedDocID {
		fmt.Println("Document id got not equal！")
		t.FailNow()
	}

	filter := bson.D{{"_id", insertedDocID}}
	update := bson.D{{"$set", bson.D{{"title", "Ledebouria socialis"}}}}
	// opts := options.Update().SetUpsert(false)
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	u, _ := uuid.Parse("de867ca7-270e-4b00-a7d6-37bd8f073058")
	identifier := []interface{}{bson.D{{"hotOptions.uuid", bson.D{{"$eq", u}}}}}
	update = bson.D{{"$set", bson.D{{"comments.$[hotOptions].content", "哈哈哈哈哈哈hahaqqqhda"}}}}
	// uuid.New().String()
	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)
	var updatedDoc models.Post
	err = collection.FindOneAndUpdate(context.TODO(), bson.D{{"_id", insertedDocID}}, update, opts).Decode(&updatedDoc)

	if err != nil {
		panic(err)
	}
	fmt.Println("=================\n", updatedDoc.Comments)

	res, _ = bson.MarshalExtJSON(updatedDoc, false, false)
	fmt.Println(string(res))

	///////////////////////////////////////////////////////////////////////
	// filter2 := bson.M{}
	statement := bson.M{"$pull": bson.M{"userActions": bson.M{"type": bson.D{{"$eq", 1}}, "userID": bson.D{{"$eq", 2}}}}}
	// result, err := collection.UpdateOne(context.TODO(), filter2, statement)
	opts = options.FindOneAndUpdate().
		// SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(context.TODO(), bson.D{{"_id", insertedDocID}}, statement, opts).Decode(&updatedDoc)

	if err != nil {
		panic(err)
	}
	fmt.Println("=================\n", updatedDoc.UserActions)

	res, _ = bson.MarshalExtJSON(updatedDoc, false, false)
	fmt.Println(string(res))

	// filter = bson.D{{"_id", insertedDocID}}
	// update = bson.D{{"$set", bson.D{{"title", "Ledebouria socialis"}}}}
	// // opts := options.Update().SetUpsert(false)
	// result, err = collection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	panic(err)
	// }

	// find := options.FindOne()
	// found := collection.FindOne(ctx, bson.D{{"_id", insertedDocID}}, find)
	// err = found.Decode(&post)
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.FailNow()
	// }
	// fmt.Println(result, "found", found)
	// fmt.Println(post.ToIndentedJSON())

}
