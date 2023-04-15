package services

type OrdersService struct{}

var ordersService *OrdersService

func GetOrdersService() *OrdersService {
	if ordersService == nil {
		ordersService = &OrdersService{}
	}
	return ordersService
}
