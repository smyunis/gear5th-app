package shared

type EntityRepository[T any] interface {
	Get(id Id) (T, error)
	Save(e T) error
}

