package api

import (
	"context"

	"github.com/Laviniarm/microservices/shipping/internal/application/core/domain"
	"github.com/Laviniarm/microservices/shipping/internal/ports"
)

type App struct{}

func New() ports.ShippingService { return &App{} }

func (a *App) Estimate(ctx context.Context, totalUnits int) (int32, error) {
	return domain.EstimateDays(totalUnits), nil
}
