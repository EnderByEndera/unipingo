package tests

import (
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewOrder(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	// 模拟admin下单
	order, err := services.GetOrdersService().NewOrder(&models.Product{
		Type:     models.MemberSubscription,
		Provider: admin.ID,
	})
	assert.Equal(t, err, nil)
	// 这里返回order是不对的，因为其ID为0。可考虑只返回创建的order的id。
	fmt.Println(order)
}

func TestGetOrder(t *testing.T) {
	order_id := "0"
	order, err := services.GetOrdersService().GetOrder(order_id)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, order, nil)
	assert.Equal(t, order.ID.String(), order_id)
	fmt.Println(order)
}

func TestPrepay(t *testing.T) {
	user_admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	product := &models.Product{
		ID:       primitive.NewObjectID(),
		Type:     models.MemberSubscription,
		Provider: user_admin.ID,
		Status:   models.ProductLaunch,
	}
	order := product.InitOrder()

	prepay_id, code, err := services.GetOrdersService().PrepayOrder(order, &user_admin)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, prepay_id, nil)
	assert.Equal(t, code, http.StatusOK)
}
