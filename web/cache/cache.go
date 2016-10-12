package cache

//Store cache store
type Store interface {
	Flush() error
	Keys() ([]string, error)
	Set(key string, val interface{}, ttl uint) error
	Get(key string, val interface{}) error
}
