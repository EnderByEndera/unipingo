package services

import (
	"context"

	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrdersService struct{}

var ordersService *OrdersService

func GetOrdersService() *OrdersService {
	if ordersService == nil {
		ordersService = &OrdersService{}
	}
	return ordersService
}

// 创建一个新订单,将关联的产品作为传入参数
func (service *OrdersService) NewOrder(product *models.Product) (order *models.Order, err error) {
	order = product.InitOrder()
	_, err = db.GetCollection("orders").InsertOne(context.Background(), order)
	if err != nil {
		return
	}
	return
}

// 返回指定状态的订单列表（"待支付"(要求将“支付失败”的放在最上面)，"已支付"等）
func (service *OrdersService) OrderList(orderStatus models.OrderStatus, page int64) (orders []*models.Order, err error) {
	filter := bson.M{}

	//当查询状态为“未支付”时，需要先查询“支付失败”的订单，两个切片拼接起来作为结果（这个if感觉相当炸裂，这样写真的合理吗！！！？？）
	if orderStatus == models.Unpaid {
		//先查询status为PaymentErro（支付失败）
		filter["status"] = models.PaymentError
		opts := options.Find().SetLimit(20).SetSkip(20 * page)
		res1, err1 := db.GetCollection("orders").Find(context.TODO(), filter, opts)
		if err1 != nil {
			return
		}
		//再查询status为Unpaid（未支付）
		filter["status"] = models.Unpaid
		opts2 := options.Find().SetLimit(20).SetSkip(20 * page)
		res2, err2 := db.GetCollection("orders").Find(context.TODO(), filter, opts2)
		if err2 != nil {
			return
		}
		//将res1和res2转为切片  后续合并
		var res11 []*models.Order
		err3 := res1.All(context.TODO(), &res11)
		if err3 != nil {
			return
		}
		var res22 []*models.Order
		err = res2.All(context.TODO(), &res22)
		//将两个切片合并
		orders = append(res11, res22...)
		return
	}
	//其他状态的查询
	if orderStatus != "" {
		filter["status"] = orderStatus
	}
	opts := options.Find().SetLimit(20).SetSkip(20 * page)
	res, err := db.GetCollection("orders").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	err = res.All(context.TODO(), &orders)
	return
}

// 修改某一个订单的状态（待支付-->已支付、待支付-->已取消、待支付-->支付失败）
func (service *OrdersService) UpdateOrderStatus(orderID primitive.ObjectID, newStatus *models.OrderStatus) (err error) {
	statement := bson.M{"$set": bson.M{"status": newStatus}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	if err != nil {
		return
	}
	res := db.GetCollection("orders").FindOneAndUpdate(context.TODO(), bson.M{"_id": orderID}, statement, opts)
	err = res.Err()
	return

}

//后续还要写的：1、支付订单（调用微信支付的API，同时要修改订单状态）2、创建产品product 3、还有啥呢？？？
