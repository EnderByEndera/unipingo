package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	SUCCESS    OrderStatus = "SUCCESS"    //支付成功
	REFUND     OrderStatus = "REFUND"     //转入退款
	NOTPAY     OrderStatus = "NOTPAY"     //未支付
	CLOSED     OrderStatus = "CLOSED"     //已关闭
	REVOKED    OrderStatus = "REVOKED"    //已撤销（付款码支付）
	USERPAYING OrderStatus = "USERPAYING" //用户支付中（付款码支付）
	PAYERROR   OrderStatus = "PAYERROR"   //支付失败(其他原因，如银行返回失败)
)

type TradeType string

const (
	JSAPI    TradeType = "JSAPI"    //公众号支付
	NATIVE   TradeType = "NATIVE"   //扫码支付
	APP      TradeType = "APP"      //APP支付
	MICROPAY TradeType = "MICROPAY" //付款码支付
	MWEB     TradeType = "MWEB"     //H5支付
	FACEPAY  TradeType = "FACEPAY"  //刷脸支付
)

type ProductStatus string

// 评论：ProductType、ProductStatus等字符串中的内容，用英文更好，和变量名一致即可。
const (
	ProductLaunch  ProductStatus = "ProductLaunch"
	ProductOffline ProductStatus = "ProductOffline"
)

type ProductType string

// 产品的类型和订单的价格value要绑定，但目前还没有确定
const (
	MemberSubscription   ProductType = "MemberSubscription"
	ModuleContentPayment ProductType = "ModuleContentPayment"
	QAPayment            ProductType = "QAPayment"
)

type TransactionAmount struct {
	Currency      *string `bson:"currency" json:"currency,omitempty"`
	PayerCurrency *string `bson:"payerCurrency" json:"payerCurrency,omitempty"`
	PayerTotal    *int64  `bson:"payerTotal" json:"payerTotal,omitempty"`
	Total         *int64  `bson:"total" json:"total,omitempty"`
}

// PromotionDetail
type PromotionDetail struct {
	CouponId            *string                `bson:"couponID" json:"couponID,omitempty"`                       // 券ID
	Name                *string                `bson:"name" json:"name,omitempty"`                               // 优惠名称
	Scope               *string                `bson:"scope" json:"scope,omitempty"`                             // GLOBAL：全场代金券；SINGLE：单品优惠
	Type                *string                `bson:"type" json:"type,omitempty"`                               // CASH：充值；NOCASH：预充值。
	Amount              *int64                 `bson:"amount" json:"amount,omitempty"`                           // 优惠券面额
	StockId             *string                `bson:"stockID" json:"stockID,omitempty"`                         // 活动ID，批次ID
	WechatpayContribute *int64                 `bson:"wechatpayContribute" json:"wechatpayContribute,omitempty"` // 单位为分
	MerchantContribute  *int64                 `bson:"merchantContribute" json:"merchantContribute,omitempty"`   // 单位为分
	OtherContribute     *int64                 `bson:"otherContribute" json:"otherContribute,omitempty"`         // 单位为分
	Currency            *string                `bson:"currency" json:"currency,omitempty"`                       // CNY：人民币，境内商户号仅支持人民币。
	GoodsDetail         []PromotionGoodsDetail `bson:"goodsDetail" json:"goodsDetail,omitempty"`
}

// PromotionGoodsDetail
type PromotionGoodsDetail struct {
	GoodsId        *string `bson:"goodsID" json:"goodsID"`                   // 商品编码
	Quantity       *int64  `bson:"quantity" json:"quantity"`                 // 商品数量
	UnitPrice      *int64  `bson:"unitPrice" json:"unitPrice"`               // 商品价格
	DiscountAmount *int64  `bson:"discountAmount" json:"discountAmount"`     // 商品优惠金额
	GoodsRemark    *string `bson:"goodsRemark" json:"goodsRemark,omitempty"` // 商品备注
}

type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionID string             `bson:"transactionID" json:"transactionID"` //微信支付订单号
	UserID        primitive.ObjectID `bson:"userID" json:"userID"`
	// 订单对应的预付单号，微信支付需要
	PrepayID         string             `bson:"prepayID" json:"prepayID"`
	TradeTypes       TradeType          `bson:"tradeType" json:"tradeType"`                       //交易类型
	Status           OrderStatus        `bson:"status" json:"status"`                             //交易状态
	TradeStateDesc   string             `bson:"tradeStateDesc" json:"tradeStateDesc"`             //交易状态描述
	BankType         string             `bson:"bankType" json:"bankType"`                         //付款银行
	Attach           string             `bson:"attach" json:"attach"`                             //附加数据
	SuccessTime      string             `bson:"successTime" json:"successTime"`                   //支付完成时间
	Amount           *TransactionAmount `bson:"amount" json:"amount,omitempty"`                   //订单金额，存储微信通知返回的金额信息
	PromotionDetails []PromotionDetail  `bson:"promotionDetail" json:"promotionDetail,omitempty"` //优惠功能
	SKUItem          SKUInfo            `bson:"skuitem" json:"skuitem"`                           //商品信息
	Value            int64              `bson:"value" json:"value"`                               //订单金额，或许在创建订单的时候用
	CreateAt         uint64             `bson:"createAt" json:"createAt"`
	CancelledAt      uint64             `bson:"cancelledAt" json:"cancelledAt"`
}

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type          ProductType        `bson:"type" json:"type"`
	Provider      primitive.ObjectID `bson:"provider" json:"provider"`
	Expiration    uint64             `bson:"expiration" json:"expiration"`
	Cut           float32            `bson:"cut" json:"cut"`
	Status        ProductStatus      `bson:"status" json:"status"`
	CreateAt      uint64             `bson:"createat" json:"createat"`
	OffsaleAt     uint64             `bson:"offsaleat" json:"offsaleat"`
	OffsaleReason string             `bson:"offsalereason" json:"offsalereason"`
}

type SKUInfo struct {
	SKUID         primitive.ObjectID `bson:"skuid" json:"skuid"`
	SKUType       ProductType        `bson:"skutype" json:"skutype"`             //产品类型
	SKUExpiration uint64             `bson:"skuexpiration" json:"skuexpiration"` //产品到期时间
	SKUStatus     ProductStatus      `bson:"skustatus" json:"skustatus"`         //产品状态
}

// 订单和用户选择的产品服务绑定，先init产品（还没有写这个函数），再使用产品来init订单
func (product *Product) InitOrder() (order *Order) {
	order = &Order{}
	order.SKUItem.SKUID = product.ID
	order.SKUItem.SKUType = product.Type
	order.SKUItem.SKUExpiration = product.Expiration
	order.SKUItem.SKUStatus = product.Status

	order.Status = NOTPAY
	//value订单的价格还没有初始化
	order.CreateAt = uint64(time.Now().Unix())
	return

}
