package server

type HttpClient interface {
	Get(url string, headers map[string]interface{}) (response HttpResponse, err error)
	Post(url string, body any, headers map[string]interface{}) (response HttpResponse, err error)
	Patch(url string, body any, headers map[string]interface{}) (response HttpResponse, err error)
	Update(url string, body any, headers map[string]interface{}) (response HttpResponse, err error)
	Delete(url string, headers map[string]interface{}) (response HttpResponse, err error)
}

type HttpResponse struct {
	StatusCode int
	Body       any
	Headers    map[string]string
}
