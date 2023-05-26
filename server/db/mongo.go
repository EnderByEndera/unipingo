package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"melodie-site/server/config"
)

type MongoConnection struct {
	Client *mongo.Client
}

var mongoConn *MongoConnection

func InitMongo() *MongoConnection {
	timeout := config.GetConfig().MongoDB.Timeout
	cli, err := mongo.NewClient(options.Client().
		ApplyURI(config.GetConfig().MongoDB.URI).
		SetTimeout(timeout))
	if err != nil {
		log.Fatal(err)
	}

	err = cli.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	mongoConnection := MongoConnection{Client: cli}

	return &mongoConnection
}

func GetMongoConn() *MongoConnection {
	if mongoConn == nil {
		mongoConn = InitMongo()
	}
	return mongoConn
}

func GetCollection(collectionName string) (collection *mongo.Collection) {
	conn := GetMongoConn()
	collection = conn.Client.Database("blog").Collection(collectionName)
	return
}

// StartSession 用来创建一个新的Session，可以支持更广泛的Abort操作和Commit操作，无需关注返回值的类型断言
func (mc *MongoConnection) StartSession(opts ...*options.SessionOptions) (session mongo.Session, err error) {
	session, err = mc.Client.StartSession(opts...)
	return
}

// UseSession 用来将已有事务函数Wrap到一个Session中以自动进行Start，Abort和Commit操作，但是返回值需要进行类型断言（如有返回值）
func (mc *MongoConnection) UseSession(opts *options.SessionOptions, transaction func(sessCtx mongo.SessionContext) (interface{}, error)) (result interface{}, err error) {
	err = mc.Client.UseSessionWithOptions(context.TODO(), opts, func(sessCtx mongo.SessionContext) error {
		err = sessCtx.StartTransaction()
		if err != nil {
			return err
		}

		// 开始执行事务
		result, err = transaction(sessCtx)

		if err != nil {
			// 事务执行失败，进行回滚
			_ = sessCtx.AbortTransaction(context.TODO())

			return err
		}

		// 提交事务
		return sessCtx.CommitTransaction(context.TODO())
	})
	return
}
