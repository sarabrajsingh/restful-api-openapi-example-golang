// server_test.go
package server_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sarabrajsingh/restful-openapi/config"
	"github.com/sarabrajsingh/restful-openapi/internal/server"
	"github.com/sarabrajsingh/restful-openapi/internal/utils"
	"github.com/sarabrajsingh/restful-openapi/mocks"
	"github.com/stretchr/testify/assert"
)

func TestIndexRedirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock dependencies
	mockLogger := mocks.NewMockLogger(ctrl)
	mockErrorStore := mocks.NewMockErrorStore(ctrl)
	bodyReader := utils.DefaultBodyReader

	// Set up mock expectations
	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()

	// Create server with mocks
	cfg, err := config.NewConfig()
	assert.NoError(t, err)
	srv := server.NewServer(cfg, mockLogger, mockErrorStore, bodyReader)

	// Create a test server
	testServer := httptest.NewServer(srv.NewRouter())
	defer testServer.Close()

	// Create a request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/", testServer.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Validate the Content-Type header
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))

	// Optionally, check the content of the body if you expect specific HTML content
	assert.Contains(t, string(body), "<!DOCTYPE html>")
	assert.Contains(t, string(body), "</html>")
}

// TestErrorsDelete tests the /errors DELETE endpoint
func TestErrorsDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock dependencies
	mockLogger := mocks.NewMockLogger(ctrl)
	mockErrorStore := mocks.NewMockErrorStore(ctrl)
	bodyReader := utils.DefaultBodyReader

	// Set up mock expectations
	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()
	mockErrorStore.EXPECT().DeleteErrors(gomock.Any()).AnyTimes()

	// Create server with mocks
	cfg, err := config.NewConfig()
	assert.NoError(t, err)
	srv := server.NewServer(cfg, mockLogger, mockErrorStore, bodyReader)

	// Create a test server
	testServer := httptest.NewServer(srv.NewRouter())
	defer testServer.Close()

	// Create a request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/errors", testServer.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestErrorsGet tests the /errors GET endpoint
func TestErrorsGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock dependencies
	mockLogger := mocks.NewMockLogger(ctrl)
	mockErrorStore := mocks.NewMockErrorStore(ctrl)
	bodyReader := utils.DefaultBodyReader

	// Set up mock expectations
	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()
	mockErrorStore.EXPECT().GetErrors(gomock.Any()).Return([]string{"error1", "error2"})

	// Create server with mocks
	cfg, err := config.NewConfig()
	assert.NoError(t, err)
	srv := server.NewServer(cfg, mockLogger, mockErrorStore, bodyReader)

	// Create a test server
	testServer := httptest.NewServer(srv.NewRouter())
	defer testServer.Close()

	// Create a request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/errors", testServer.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "error1")
	assert.Contains(t, string(body), "error2")
}
