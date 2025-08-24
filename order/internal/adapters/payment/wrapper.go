package paymentadapter

import (
	"context"

	"github.com/Laviniarm/microservices/order/internal/application/core/domain"
)

type ClientWrapper struct {
	inner *Adapter // teu Adapter já existente (que tem Charge)
}

func NewClientWrapper(inner *Adapter) *ClientWrapper { return &ClientWrapper{inner: inner} }

func (w *ClientWrapper) Process(ctx context.Context, req struct {
	OrderID    int64
	CustomerID int64
	Items      []domain.OrderItem
	Total      float32
	Currency   string
}) (bool, error) {
	// monta uma Order mínima pro Charge
	o := domain.NewOrder(req.CustomerID, req.Items) // <- retorna 'domain.Order' (valor)
	o.ID = req.OrderID                              // define o ID vindo da requisição

	if err := w.inner.Charge(&o); err != nil { // <- passe o endereço &o
		return false, err
	}
	return true, nil
}
