package http

type Headers map[string]string
type Params map[string]string

type HTTPService interface {
	Get(url string, headers Headers, params Params) ([]byte, error)
	Post(url string, headers Headers, data interface{}) ([]byte, error)
	Put(url string, headers Headers, data interface{}) ([]byte, error)
	Delete(url string, headers Headers, data interface{}) ([]byte, error)
}
