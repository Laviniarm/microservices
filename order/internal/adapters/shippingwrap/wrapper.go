package shippingwrap

import (
	"context"

	scli "github.com/Laviniarm/microservices/order/internal/adapters/shipping"
)

type ClientWrapper struct {
	inner *scli.Adapter
}

func NewClientWrapper(inner *scli.Adapter) *ClientWrapper { return &ClientWrapper{inner: inner} }

func (w *ClientWrapper) Estimate(ctx context.Context, orderID int64, items []struct {
	ID  string
	Qty int32
}) (int32, error) {
	payload := make([]scli.Item, 0, len(items))
	for _, it := range items {
		payload = append(payload, scli.Item{ID: it.ID, Qty: it.Qty})
	}
	return w.inner.Estimate(ctx, orderID, payload)
}
