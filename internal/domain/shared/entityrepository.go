package shared

type EntityRepository[T any] interface {
	Get(id ID) (T, error)
	Save(e T) error
}
