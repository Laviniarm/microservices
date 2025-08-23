package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Laviniarm/microservices/order/internal/ports"
)

type InventoryRepo struct{ db *sql.DB }

func NewInventoryRepo(db *sql.DB) *InventoryRepo { return &InventoryRepo{db: db} }

var _ ports.InventoryRepository = (*InventoryRepo)(nil)

func (r *InventoryRepo) MissingIDs(ctx context.Context, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		ph[i] = "?"
		args[i] = id
	}

	q := fmt.Sprintf(`SELECT id FROM inventory_items WHERE id IN (%s)`, strings.Join(ph, ","))
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := map[string]struct{}{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		found[id] = struct{}{}
	}
	var missing []string
	for _, id := range ids {
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing, rows.Err()
}
