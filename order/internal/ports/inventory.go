package ports

import "context"

type InventoryRepository interface {
	MissingIDs(ctx context.Context, ids []string) ([]string, error)
}
