package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sarabrajsingh/restful-openapi/internal/handlers"
	"github.com/sarabrajsingh/restful-openapi/internal/logging"
	"github.com/sarabrajsingh/restful-openapi/internal/models"
	"github.com/sarabrajsingh/restful-openapi/internal/utils"
	"github.com/sarabrajsingh/restful-openapi/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)

	testCases := []struct {
		description      string
		mockGetErrors    func(logging.Logger) []string
		expectedStatus   int
		expectedResponse models.GetErrorsResponse
	}{
		{
			description: "No errors",
			mockGetErrors: func(log logging.Logger) []string {
				return []string{}
			},
			expectedStatus: http.StatusOK,
			expectedResponse: models.GetErrorsResponse{
				Errors: []string{},
			},
		},
		{
			description: "One error",
			mockGetErrors: func(log logging.Logger) []string {
				return []string{"error1"}
			},
			expectedStatus: http.StatusOK,
			expectedResponse: models.GetErrorsResponse{
				Errors: []string{"error1"},
			},
		},
		{
			description: "Multiple errors",
			mockGetErrors: func(log logging.Logger) []string {
				return []string{"error1", "error2"}
			},
			expectedStatus: http.StatusOK,
			expectedResponse: models.GetErrorsResponse{
				Errors: []string{"error1", "error2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := handlers.GetErrors(mockLogger, tc.mockGetErrors)

			req, err := http.NewRequest("GET", "/api/v1/errors", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if status := w.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			var actualResponse models.GetErrorsResponse
			if err := json.NewDecoder(w.Body).Decode(&actualResponse); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if len(actualResponse.Errors) != len(tc.expectedResponse.Errors) {
				t.Errorf("Expected Errors: %v, got: %v", tc.expectedResponse.Errors, actualResponse.Errors)
			}

			for i, errorStr := range actualResponse.Errors {
				if errorStr != tc.expectedResponse.Errors[i] {
					t.Errorf("Expected Error at index %d: %v, got: %v", i, tc.expectedResponse.Errors[i], errorStr)
				}
			}
		})
	}
}

func TestDeleteErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)

	testCases := []struct {
		description      string
		mockDeleteErrors func(logging.Logger)
		expectedStatus   int
	}{
		{
			description:      "Successful deletion of errors",
			mockDeleteErrors: func(log logging.Logger) {},
			expectedStatus:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := handlers.DeleteErrors(mockLogger, tc.mockDeleteErrors)

			req, err := http.NewRequest("DELETE", "/api/v1/errors", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if status := w.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}

func TestTempPostHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		description      string
		requestBody      string
		mockAddErrorFunc func(logging.Logger, string)
		expectedStatus   int
		expectedResponse models.TempPostResponse
	}{
		{
			description:      "Valid request; Overtemp",
			requestBody:      `{"data":"1234:1721964434:'Temperature':95.0"}`,
			mockAddErrorFunc: func(log logging.Logger, data string) {},
			expectedStatus:   http.StatusOK,
			expectedResponse: models.TempPostResponse{
				DeviceId:      1234,
				Overtemp:      true,
				FormattedTime: "2024/07/25 23:27:14",
			},
		},
		{
			description:      "Valid request; Not Overtemp",
			requestBody:      `{"data":"1234:1721964434:'Temperature':89.9"}`,
			mockAddErrorFunc: func(log logging.Logger, data string) {},
			expectedStatus:   http.StatusOK,
			expectedResponse: models.TempPostResponse{
				Overtemp: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockLogger := mocks.NewMockLogger(ctrl)
			mockLogger.EXPECT().Println(gomock.Any()).AnyTimes()

			bodyReader := func(r io.Reader) ([]byte, error) {
				return io.ReadAll(r)
			}

			handler := handlers.TempPost(mockLogger, tc.mockAddErrorFunc, bodyReader)

			req, err := http.NewRequest("POST", "/api/v1/temp", strings.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if status := w.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			var actualResponse models.TempPostResponse

			if err := json.NewDecoder(w.Body).Decode(&actualResponse); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			assert.NotNil(t, actualResponse)
			assert.Equal(t, actualResponse, tc.expectedResponse)
		})
	}
}

func TestTempPostUnhappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		description      string
		requestBody      string
		mockAddErrorFunc func(logging.Logger, string)
		expectedStatus   int
		expectedResponse models.Response400
		bodyReader       utils.BodyReaderFunc
	}{
		{
			description:      "Bad Request; malformed JSON request payload",
			mockAddErrorFunc: func(log logging.Logger, data string) {},
			expectedStatus:   http.StatusBadRequest,
			bodyReader: func(io.Reader) ([]byte, error) {
				return nil, errors.New("foobar error")
			},
			expectedResponse: models.Response400{
				Error: "Failed to parse request body",
			},
		},
		{
			description:      "Bad Request; unable to unmarshal JSON payload",
			requestBody:      `{"data":}`,
			mockAddErrorFunc: func(log logging.Logger, data string) {},
			expectedStatus:   http.StatusBadRequest,
			bodyReader: func(r io.Reader) ([]byte, error) {
				return io.ReadAll(r)
			},
			expectedResponse: models.Response400{
				Error: "Failed to parse JSON",
			},
		},
		{
			description:      "Bad Request; invalid fields in JSON payload",
			requestBody:      `{"data":"abc:def:'Temperature':95.0"}`,
			mockAddErrorFunc: func(log logging.Logger, data string) {},
			expectedStatus:   http.StatusBadRequest,
			bodyReader: func(r io.Reader) ([]byte, error) {
				return io.ReadAll(r)
			},
			expectedResponse: models.Response400{
				Error: "bad request",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockLogger := mocks.NewMockLogger(ctrl)
			mockLogger.EXPECT().Println(gomock.Any()).AnyTimes()

			handler := handlers.TempPost(mockLogger, tc.mockAddErrorFunc, tc.bodyReader)

			req, err := http.NewRequest("POST", "/api/v1/temp", strings.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if status := w.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			var actualResponse models.Response400

			if err := json.NewDecoder(w.Body).Decode(&actualResponse); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			assert.NotNil(t, actualResponse)
			assert.Equal(t, actualResponse, tc.expectedResponse)
		})
	}
}
