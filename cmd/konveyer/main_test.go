package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetIndex(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	body := bytes.NewBuffer([]byte("{}"))

	req, err := http.NewRequest("GET", "/api/v1/", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Get Index failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/ failed with error code %d.", resp.Code)
	}
}
