package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type TestRequest struct {
	bearerToken string
	Router      *gin.Engine
}

func (r TestRequest) makeRequest(method, url string, body any) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if r.bearerToken != "" {
		req.Header.Add("Authorization", r.bearerToken)
	}

	if method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	writer := httptest.NewRecorder()

	r.Router.ServeHTTP(writer, req)

	return writer
}

func (r *TestRequest) SetBearerToken(token string) {
	r.bearerToken = token
}

func (r TestRequest) Post(url string, body any) *httptest.ResponseRecorder {
	return r.makeRequest("POST", url, body)
}

func (r TestRequest) Put(url string, body any) *httptest.ResponseRecorder {
	return r.makeRequest("PUT", url, body)
}
func (r TestRequest) Patch(url string, body any) *httptest.ResponseRecorder {
	return r.makeRequest("PATCH", url, body)
}
func (r TestRequest) Get(url string, body any) *httptest.ResponseRecorder {
	return r.makeRequest("GET", url, body)
}
func (r TestRequest) Delete(url string, body any) *httptest.ResponseRecorder {
	return r.makeRequest("DELETE", url, body)
}
