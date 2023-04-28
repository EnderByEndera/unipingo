package services

import (
	"context"

	"melodie-site/server/config"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type OrdersService struct{}

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
func (service *OrdersService) NewOrder(product *models.Product) (order *models.Order, err error) {
	order = product.InitOrder()
	_, err = db.GetCollection("orders").InsertOne(context.Background(), order)
	if err != nil {
		return
	}
	return
}

func (service *OrdersService) PrepayOrder(order *models.Order, user *models.User) (prepayID *string, statusCode int, err error) {
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

	// 微信支付会返回状态码
	statusCode = res.Response.StatusCode
	if err != nil {
		/*
			switch {
			case core.IsAPIError(err, "SYSTEM_ERROR"):
				err = errors.New("系统异常，请用相同参数重新调用")
			case core.IsAPIError(err, "SIGN_ERROR"):
				err = errors.New("请检查签名参数和方法是否都符合签名算法要求")
			case core.IsAPIError(err, "RULE_LIMIT"):
				err = errors.New("因业务规则限制请求频率")
			case core.IsAPIError(err, "PARAM_ERROR"):
				err = errors.New("请检查请求参数")
			case core.IsAPIError(err, "OUT_TRADE_NO_USED"):
				err = errors.New("请核实商户订单号是否重复提交")
			case core.IsAPIError(err, "ORDER_NOT_EXIST"):
				err = errors.New("请检查订单是否发起过交易")
			case core.IsAPIError(err, "ORDER_CLOSED"):
				err = errors.New("当前订单已关闭，请重新下单")
			case core.IsAPIError(err, "OPENID_MISMATCH"):
				err = errors.New("请确认openid和appid是否匹配")
			case core.IsAPIError(err, "NO_AUTH"):
				err = errors.New("请商户前往申请此接口相关权限")
			}
		*/
		return
	}
	prepayID = resp.PrepayId
	updatePrepay := bson.D{{Key: "$set", Value: bson.M{"prepayID": *prepayID}}}
	err = db.GetCollection("orders").FindOneAndUpdate(context.TODO(), bson.M{"_id": order.ID}, updatePrepay).Err()
	return
}

func (service *OrdersService) GetOrder(order_id string) (order *models.Order, err error) {
	filter := bson.M{"_id": order_id}
	order = &models.Order{}
	err = db.GetCollection("orders").FindOne(context.TODO(), filter).Decode(order)
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
