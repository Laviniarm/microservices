package db

import (
	"fmt"

	"github.com/Laviniarm/microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order
	res := a.db.Preload("OrderItems").First(&orderEntity, id)
	if res.Error != nil {
		return domain.Order{}, res.Error
	}

	var orderItems []domain.OrderItem
	for _, oi := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: oi.ProductCode,
			UnitPrice:   oi.UnitPrice,
			Quantity:    oi.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.Unix(),
	}

	return order, nil
}

func (a *Adapter) Save(order *domain.Order) error {
	var orderItems []OrderItem
	for _, oi := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: oi.ProductCode,
			UnitPrice:   oi.UnitPrice,
			Quantity:    oi.Quantity,
		})
	}

	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,
	}

	res := a.db.Create(&orderModel)
	if res.Error == nil {
		order.ID = int64(orderModel.ID)
		order.CreatedAt = orderModel.CreatedAt.Unix()
	}

	return res.Error
}

func (a *Adapter) UpdateStatus(orderID int64, status string) error {
	res := a.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status)
	return res.Error
}
