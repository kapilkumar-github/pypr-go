package variable

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kapilkumar9395/pypr/internal/infra/cache"
)

type VariableRepo struct {
	db    *pgxpool.Pool
	cache *cache.Cache[uuid.UUID]
}

func NewVariableRepo(db *pgxpool.Pool, cache *cache.Cache[uuid.UUID]) *VariableRepo {
	return &VariableRepo{db, cache}
}

func (r *VariableRepo) GetIDsByCodes(
	ctx context.Context,
	codes []string,
) (map[string]uuid.UUID, error) {

	result := make(map[string]uuid.UUID)
	missing := make([]string, 0)

	// 1. Check cache
	for _, code := range codes {
		if id, ok := r.cache.Get(code); ok {
			result[code] = id
		} else {
			missing = append(missing, code)
		}
	}

	// Everything was cached
	if len(missing) == 0 {
		return result, nil
	}

	// 2. Query DB only for missing codes
	rows, err := r.db.Query(ctx, `
        SELECT id, key
        FROM variables
        WHERE key = ANY($1)
    `, missing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   uuid.UUID
			code string
		)

		if err := rows.Scan(&id, &code); err != nil {
			return nil, err
		}

		result[code] = id
		r.cache.Set(code, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
