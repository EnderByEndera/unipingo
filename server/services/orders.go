package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"melodie-site/server/config"
	"melodie-site/server/db"
	"melodie-site/server/models"
	myutils "melodie-site/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type OrdersService struct {
	userAndOrderMutex *myutils.KeyedMutex
}

type UserOrderInterface struct {
	UserID  primitive.ObjectID
	OrderID primitive.ObjectID
}

var (
	ordersService *OrdersService
	wechatClient  *core.Client

	// TODO: 目前需要等待产品注册微信商户后获得以下三个参数
	mchID                      string = config.GetConfig().WECHAT.MCHID            // 商户号
	mchCertificateSerialNumber string = config.GetConfig().WECHAT.MCHCERTSERIALNUM // 商户证书序列号
	mchAPIv3Key                string = config.GetConfig().WECHAT.MCHAPIV3KEY      // 商户APIv3密钥
)

func GetOrdersService() *OrdersService {
	if ordersService == nil {
		ordersService = &OrdersService{}
	}
	return ordersService
}

func (ordersService *OrdersService) LockUserAndOrder(userID, orderID primitive.ObjectID) {
	key := UserOrderInterface{UserID: userID, OrderID: orderID}
	ordersService.userAndOrderMutex.Lock(key)
}

func (ordersService *OrdersService) UnlockUserAndOrder(userID, ansID primitive.ObjectID) {
	key := IDInterface{UserID: userID, AnswerID: ansID}
	ordersService.userAndOrderMutex.Unlock(key)
}

func GetWechatClient() (client *core.Client, err error) {
	if wechatClient != nil {
		client = wechatClient
		return
	}

	// 加载商户API证书私钥
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("cert/apiclient_key.pem")
	if err != nil {
		return
	}

	ctx := context.Background()
	// 初始化wechatClient所需的参数
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID,
			mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	// 建立wechatClient
	client, err = core.NewClient(ctx, opts...)
	wechatClient = client
	return
}

// 创建一个新订单,将关联的产品作为传入参数
func (service *OrdersService) NewOrder(userID primitive.ObjectID, product *models.Product) (order *models.Order, err error) {
	order = product.InitOrder()
	order.UserID = userID
	_, err = db.GetCollection("orders").InsertOne(context.Background(), order)
	if err != nil {
		return
	}
	return
}

func (service *OrdersService) PrepayOrder(order *models.Order, user *models.User) (prepayID string, res *core.APIResult, err error) {
	appid := config.GetConfig().WECHAT.APPID
	client, err := GetWechatClient()
	if err != nil {
		return
	}

	jsapiService := jsapi.JsapiApiService{Client: client}
	resp, res, err := jsapiService.Prepay(context.TODO(), jsapi.PrepayRequest{
		Appid:       core.String(appid),
		Mchid:       core.String(mchID),
		Description: core.String(string(order.SKUItem.SKUType)),
		OutTradeNo:  core.String(order.ID.String()),
		Attach:      core.String(string(order.Status)),
		//  微信支付建议订单有效期为5分钟
		TimeExpire: core.Time(time.UnixMicro(int64(order.CreateAt)).Add(5 * time.Minute)),

		// TODO: 需要进一步沟通具体API名称
		// 回调URL，用以之后微信支付服务端异步通知后端更新订单状态
		NotifyUrl: core.String("https://9325061_wechatapi.houzhanyi.com/api/orders/notify"),

		Amount: &jsapi.Amount{
			Total:    core.Int64(order.Value),
			Currency: core.String("CNY"), // 可以不必要，但是为了代码可读性故填上
		},
		Payer: &jsapi.Payer{
			Openid: core.String(user.WechatInfo.OpenID),
		},
	})

	if err != nil {
		return
	}

	prepayID = *resp.PrepayId
	updatePrepay := bson.D{{Key: "$set", Value: bson.M{"prepayID": prepayID}}}
	err = db.GetCollection("orders").FindOneAndUpdate(context.TODO(), bson.M{"_id": order.ID}, updatePrepay).Err()
	return
}

func (service *OrdersService) NotifyOrder(user *models.User, request *http.Request) (statusCode int, error error) {
	userID := user.ID
	_, err := GetWechatClient()
	if err != nil {
		return
	}
	//1. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	//2. 初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	//将解密后的内容封装在Transaction中
	transaction := new(payments.Transaction)
	//3、验签+解密一体
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		statusCode = http.StatusUnauthorized
		return
	}

	// 处理通知内容
	orderID, err := primitive.ObjectIDFromHex(*transaction.OutTradeNo)
	if err != nil {
		statusCode = http.DefaultMaxHeaderBytes
		return
	}
	orderStatus := (*models.OrderStatus)(transaction.TradeState)
	// 1、加锁
	ordersService.LockUserAndOrder(userID, orderID)
	defer ordersService.UnlockUserAndOrder(userID, orderID)

	// 2、先查询数据库，防止收集重复信息
	order, err := ordersService.GetOrder(orderID)
	if err != nil {
		statusCode = http.StatusNotFound
		return
	}
	//数据库里面订单的状态
	orderOldStatus := order.Status
	//满足所有一个条件就返回,表明是重复信息或者过期通知

	if (*orderStatus == models.NOTPAY) ||
		(orderOldStatus == *orderStatus) ||
		(*orderStatus == models.SUCCESS && orderOldStatus == models.REFUND) ||
		((*orderStatus == models.USERPAYING || *orderStatus == models.PAYERROR || *orderStatus == models.REVOKED) && orderOldStatus == models.CLOSED) {
		//TODO  这里状态码返回多少呢，如果不是200 还要设置错误
		return
	}

	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)

	// 3、这里要更新order的很多信息   数据库里面对象是order  微信支付通知对象是transcation
	transactionID, err := primitive.ObjectIDFromHex(*transaction.TransactionId)
	if err != nil {
		statusCode = http.StatusBadRequest
		return
	}
	order.TransactionID = transactionID
	order.TradeTypes = models.TradeType(*transaction.TradeType)
	order.Status = *orderStatus
	order.TradeStateDesc = *transaction.TradeStateDesc
	order.BankType = *transaction.BankType
	order.Attach = *transaction.Attach
	order.SuccessTime = *transaction.SuccessTime
	openID, err := primitive.ObjectIDFromHex(*transaction.Payer.Openid)
	if err != nil {
		statusCode = http.StatusBadRequest
		return
	}
	order.UserID = openID
	order.Amount = (*models.TransactionAmount)(transaction.Clone().Amount)
	//在转换promotionDetail的时候用上面的方式不能强转，这里就先用了json反序列化
	pd := (transaction.Clone().PromotionDetail)
	pd1, err := json.Marshal(&pd)
	if err != nil {
		statusCode = http.StatusBadRequest
		return
	}
	json.Unmarshal(pd1, &order.PromotionDetails)

	//4、更新到数据库
	statement := bson.M{"$set": &order}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	res := db.GetCollection("order").FindOneAndUpdate(context.TODO(), bson.M{"_id": orderID}, statement, opts)
	err = res.Err()
	if err != nil {
		statusCode = http.StatusNotFound
		return
	}
	statusCode = http.StatusOK
	return

}

func (service *OrdersService) GetOrder(order_id primitive.ObjectID) (order *models.Order, err error) {
	filter := bson.M{"_id": order_id}
	order = &models.Order{}
	err = db.GetCollection("orders").FindOne(context.TODO(), filter).Decode(order)
	return
}

// 返回指定状态的订单列表（"待支付"(要求将“支付失败”的放在最上面)，"已支付"等）
func (service *OrdersService) OrderList(orderStatus models.OrderStatus, page int64) (orders []*models.Order, err error) {
	filter := bson.M{}

	//当查询状态为“未支付”时，需要先查询“支付失败”的订单，两个切片拼接起来作为结果（这个if感觉相当炸裂，这样写真的合理吗！！！？？）
	if orderStatus == models.NOTPAY {
		//先查询status为PaymentErro（支付失败）
		filter["status"] = models.PAYERROR
		opts := options.Find().SetLimit(20).SetSkip(20 * page)
		res1, err1 := db.GetCollection("orders").Find(context.TODO(), filter, opts)
		if err1 != nil {
			return
		}
		//再查询status为Unpaid（未支付）
		filter["status"] = models.NOTPAY
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

func (service *OrdersService) GetOrderStatus(order *models.Order) (orderStatus models.OrderStatus, res *core.APIResult, err error) {
	if err != nil {
		return
	}
	client, err := GetWechatClient()
	if err != nil {
		return
	}
	jsapiService := jsapi.JsapiApiService{Client: client}
	resp, res, err := jsapiService.QueryOrderById(context.Background(), jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(order.TransactionID.String()),
		Mchid:         core.String(mchID),
	})
	if err != nil {
		return
	}
	orderStatus = models.OrderStatus(*resp.TradeState)
	return
}
