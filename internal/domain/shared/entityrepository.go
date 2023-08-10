package shared

import "context"

type EntityRepository[T any] interface {
	Get(ctx context.Context, id ID) (T, error)
	Save(ctx context.Context, e T) error
}
