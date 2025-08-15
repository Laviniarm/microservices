package api

import (
	"github.com/Laviniarm/microservices/order/internal/application/core/domain"
	"github.com/Laviniarm/microservices/order/internal/ports"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := a.db.Save(&order)
	if err != nil {
		order.Status = "cancelled"
		_ = a.db.UpdateStatus(order.ID, order.Status)
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}
	order.Status = "Paid"
	_ = a.db.UpdateStatus(order.ID, order.Status)

	return order, nil
}
