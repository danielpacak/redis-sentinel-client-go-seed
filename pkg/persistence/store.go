package persistence

type Store interface {
	Set(key, value string) error
	Keys() ([]string, error)
	Info() ([]string, error)
}
