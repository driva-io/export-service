package server

import (
	"errors"
	"strings"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]interface{}
	Body    any
}

type MockHttpClient struct {
	requests  []Request
	responses map[string]HttpResponse
}

func NewMockHttpClient() *MockHttpClient {
	return &MockHttpClient{
		responses: make(map[string]HttpResponse),
	}
}

func (m *MockHttpClient) Expect(method, url string, response HttpResponse) {
	m.responses[method+" "+url] = response
}

func (m *MockHttpClient) VerifyRequest(method, url string) (Request, bool) {
	for _, r := range m.requests {
		if r.Method == method && strings.Contains(r.URL, url) {
			return r, true
		}
	}
	return Request{}, false
}

func (m *MockHttpClient) handle(method, url string, headers map[string]interface{}, body any) (HttpResponse, error) {
	m.requests = append(m.requests, Request{Method: method, URL: url, Headers: headers, Body: body})
	for k, v := range m.responses {
		expectedRoute := strings.ReplaceAll(k, method+" ", "")
		if strings.HasPrefix(k, method) && strings.Contains(url, expectedRoute) {
			return v, nil
		}
	}
	return HttpResponse{}, errors.New("HTTP MOCK unexpected request: " + method + " " + url)
}

func (m *MockHttpClient) Get(url string, headers map[string]interface{}) (response HttpResponse, err error) {
	return m.handle("GET", url, headers, nil)
}

func (m *MockHttpClient) Post(url string, body any, headers map[string]interface{}) (response HttpResponse, err error) {
	return m.handle("POST", url, headers, body)
}

func (m *MockHttpClient) Patch(url string, body any, headers map[string]interface{}) (response HttpResponse, err error) {
	return m.handle("PATCH", url, headers, body)
}

func (m *MockHttpClient) Update(url string, body any, headers map[string]interface{}) (response HttpResponse, err error) {
	return m.handle("PATCH", url, headers, body)
}

func (m *MockHttpClient) Delete(url string, headers map[string]interface{}) (response HttpResponse, err error) {
	return m.handle("DELETE", url, headers, nil)
}
