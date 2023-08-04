package shared

type EntityRepository[T any] interface {
	get(id Id) (T, error)
	save(e T) error
}

