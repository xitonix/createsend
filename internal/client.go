package internal

type Client interface {
	Get(path string, result interface{}) error
	Post(path string, result, body interface{}) error
	Put(path string, result, body interface{}) error
	Delete(path string) error
}
