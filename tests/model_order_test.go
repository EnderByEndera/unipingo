package tests

import (
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrder(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	// 模拟admin下单
	order, err := services.GetOrdersService().NewOrder(admin.ID,&models.Product{
		Type:     models.MemberSubscription,
		Provider: admin.ID,
	})
	assert.Equal(t, err, nil)
	// 这里返回order是不对的，因为其ID为0。可考虑只返回创建的order的id。
	fmt.Println(order)
}


func TestUpdateOrderStatus(t *testing.T){
	orderID,_ := primitive.ObjectIDFromHex("6443fb6d200eb1117c4094b4")
	orderStatus := models.PaymentError
	err := services.GetOrdersService().UpdateOrderStatus(orderID, &orderStatus)
	assert.Equal(t, err, nil)
}

func TestListOrderByStatus(t *testing.T){
	orderStatus := models.Paid
	page:=0
	orders, err := services.GetOrdersService().OrderList(orderStatus, int64(page))
	assert.Equal(t, err, nil)
	for _,order :=range orders{
		fmt.Println(order.ID)
	}

}
