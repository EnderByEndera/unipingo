package routers

import (
	"context"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PrepayOrder(c *gin.Context) {
	// TODO: 具体的OPENID获取形式需要和前端进行沟通，前端不一定直接传输ID过来
	open_id := c.Request.Form.Get("openid")
	var user *models.User

	if open_id == "" {
		// 尝试从JWTToken中获得
		claims, err := utils.GetClaims(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		user, err = services.GetAuthService().GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
	} else {
		var err error
		user, err = services.GetAuthService().GetUserByWechatOpenID(open_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	// TODO: 获取订单的形式需要和前端沟通
	order_id := c.Request.Form.Get("out_trade_no")
	orderID, err := primitive.ObjectIDFromHex(order_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	order, err := services.GetOrdersService().GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	prepay_id, res, err := services.GetOrdersService().PrepayOrder(order, user)
	if err != nil {
		c.JSON(res.Response.StatusCode, err.Error())
	}

	newOrderStatus := models.NOTPAY
	services.GetOrdersService().UpdateOrderStatus(order.ID, &newOrderStatus)

	// TODO: 暂时确定只返回prepay_id，如需其他数据再进行修改
	c.JSON(http.StatusOK, gin.H{
		"prepay_id": prepay_id,
	})
}
func NotifyOrder(c *gin.Context) {
	//1. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(services.MchID)
	//2. 初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(services.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	//将解密后的内容封装在Transaction中
	transaction := new(payments.Transaction)
	//3、验签+解密一体
	_, err := handler.ParseNotifyRequest(context.Background(), c.Request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.GetOrdersService().NotifyOrder(transaction)

	//接收成功：HTTP应答状态码需返回200或204，无需返回应答报文。
	//接收失败：HTTP应答状态码需返回5XX或4XX，同时需返回应答报文
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
	} else {
		c.JSON(http.StatusOK, nil)
	}
}
