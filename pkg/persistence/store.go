package persistence

type Store interface {
	Keys() ([]string, error)
}
