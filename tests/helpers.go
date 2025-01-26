package tests

import (
	"net/http"
	"polling_websocket/pkg/domain/models"
	"time"
)

type MockHTTPClient struct {
	DoFunc        func(req *http.Request) (*http.Response, error)
	DoRequestFunc func(method, url, authToken string, body interface{}) ([]byte, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (m *MockHTTPClient) DoRequest(method, url, authToken string, body interface{}) ([]byte, error) {
	return m.DoRequestFunc(method, url, authToken, body)
}

func stringPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64    { return &i }
func boolPtr(b bool) *bool       { return &b }
func float64Ptr(f float64) *float64 {
	return &f
}

func customTime(t time.Time) *models.CustomTime {
	return &models.CustomTime{Time: t}
}
