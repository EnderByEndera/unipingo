package tests

import (
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestOrder(t *testing.T) {
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
