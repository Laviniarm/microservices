package ports

import "context"

type ShippingService interface {
	Estimate(ctx context.Context, totalUnits int) (int32, error)
}
