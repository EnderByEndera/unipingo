package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	Unpaid       OrderStatus = "待支付"
	Paid         OrderStatus = "已支付"
	Cancelled    OrderStatus = "已取消"
	PaymentError OrderStatus = "支付失败" //也属于“待支付”，但是为了满足将“支付失败”的订单展示在最前面的要求，所以这里单独列出来
)

type ProductStatus string

// 评论：ProductType、ProductStatus等字符串中的内容，用英文更好，和变量名一致即可。
const (
	ProductLaunch  ProductStatus = "产品上线"
	ProductOffline ProductStatus = "产品下线"
)

type ProductType string

// 产品的类型和订单的价格value要绑定，但目前还没有确定
const (
	MemberSubscription   ProductType = "会员订阅"
	ModuleContentPayment ProductType = "模块内容付费"
	QAPayment            ProductType = "问答付费"
)

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"userID" json:"userID"`
	// 订单对应的预付单号，微信支付需要
	PrepayID    string      `bson:"prepayID" json:"prepayID"`
	SKUItem     SKUInfo     `bson:"skuitem" json:"skuitem"`
	Status      OrderStatus `bson:"status" json:"status"`
	Value       int64       `bson:"value" json:"value"`
	CreateAt    uint64      `bson:"createAt" json:"createAt"`
	CancelledAt uint64      `bson:"cancelledAt" json:"cancelledAt"`
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
	SKUType       ProductType        `bson:"skutype" json:"skutype"`
	SKUExpiration uint64             `bson:"skuexpiration" json:"skuexpiration"`
	SKUStatus     ProductStatus      `bson:"skustatus" json:"skustatus"`
}

// 订单和用户选择的产品服务绑定，先init产品（还没有写这个函数），再使用产品来init订单
func (product *Product) InitOrder() (order *Order) {
	order = &Order{}
	order.SKUItem.SKUID = product.ID
	order.SKUItem.SKUType = product.Type
	order.SKUItem.SKUExpiration = product.Expiration
	order.SKUItem.SKUStatus = product.Status

	order.UserID = product.Provider // 没有理解这段代码的含义
	order.Status = Unpaid
	//value订单的价格还没有初始化
	order.CreateAt = uint64(time.Now().UnixMicro()) // 这里建议改为Micro更好
	return

}
