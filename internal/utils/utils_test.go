package utils_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sarabrajsingh/restful-openapi/internal/models"
	"github.com/sarabrajsingh/restful-openapi/internal/utils"
)

func TestPayloadParserHelper(t *testing.T) {
	testCases := []struct {
		description     string
		payloadData     string
		expectedPayload *models.TempPostPayload
		expectedError   string
	}{
		{
			description:     "Valid payload",
			payloadData:     "1234:1721964434:'Temperature':95.0",
			expectedPayload: &models.TempPostPayload{DeviceId: 1234, EpochMS: 1721964434, Temperature: 95.0},
			expectedError:   "",
		},
		{
			description:     "Invalid payload - incorrect number of arguments",
			payloadData:     "1234:1721964434:95.0",
			expectedPayload: nil,
			expectedError:   "invalid number of arguments in request body",
		},
		{
			description:     "Invalid payload - device ID is not an int",
			payloadData:     "abc:1721964434:'Temperature':95.0",
			expectedPayload: nil,
			expectedError:   "could not parse device_id=abc to an int32",
		},
		{
			description:     "Invalid payload - epochMS is not an int",
			payloadData:     "1234:abc:'Temperature':95.0",
			expectedPayload: nil,
			expectedError:   "could not parse epochMS=abc to an int64",
		},
		{
			description:     "Invalid payload - temperature key is mislabeled",
			payloadData:     "1234:1721964434:'Temp':95.0",
			expectedPayload: nil,
			expectedError:   "temperature key is mislabelled: 'Temp'",
		},
		{
			description:     "Invalid payload - temperature is not a float",
			payloadData:     "1234:1721964434:'Temperature':abc",
			expectedPayload: nil,
			expectedError:   "could not parse temperature=abc to a float64",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualPayload, actualError := utils.PayloadParserHelper(tc.payloadData)

			if actualError != nil && tc.expectedError == "" {
				t.Errorf("Unexpected error: %v", actualError)
			}

			if actualError == nil && tc.expectedError != "" {
				t.Errorf("Expected error: %v, got nil", tc.expectedError)
			}

			if actualError != nil && actualError.Error() != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, actualError.Error())
			}

			if actualPayload != nil && tc.expectedPayload == nil {
				t.Errorf("Expected nil payload, got: %v", actualPayload)
			}

			if actualPayload == nil && tc.expectedPayload != nil {
				t.Errorf("Expected payload: %v, got nil", tc.expectedPayload)
			}

			if actualPayload != nil && tc.expectedPayload != nil {
				if *actualPayload != *tc.expectedPayload {
					t.Errorf("Expected payload: %v, got: %v", tc.expectedPayload, actualPayload)
				}
			}
		})
	}
}

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name                 string
		errorMessage         string
		statusCode           int
		expectedContentType  string
		expectedStatusCode   int
		expectedResponseBody models.Response400
	}{
		{
			name:                 "Valid Error Message & Status",
			errorMessage:         "Invalid request",
			statusCode:           400,
			expectedContentType:  "application/json; charset=UTF-8",
			expectedStatusCode:   400,
			expectedResponseBody: models.Response400{Error: "Invalid request"},
		},
		{
			name:                 "Server Error",
			errorMessage:         "Internal Server Error",
			statusCode:           500,
			expectedContentType:  "application/json; charset=UTF-8",
			expectedStatusCode:   500,
			expectedResponseBody: models.Response400{Error: "Internal Server Error"},
		},
		{
			name:                 "Unauthorized Access",
			errorMessage:         "Unauthorized",
			statusCode:           401,
			expectedContentType:  "application/json; charset=UTF-8",
			expectedStatusCode:   401,
			expectedResponseBody: models.Response400{Error: "Unauthorized"},
		},
		{
			name:                 "Not Found",
			errorMessage:         "Resource not found",
			statusCode:           404,
			expectedContentType:  "application/json; charset=UTF-8",
			expectedStatusCode:   404,
			expectedResponseBody: models.Response400{Error: "Resource not found"},
		},
		{
			name:                 "Forbidden Access",
			errorMessage:         "Forbidden",
			statusCode:           403,
			expectedContentType:  "application/json; charset=UTF-8",
			expectedStatusCode:   403,
			expectedResponseBody: models.Response400{Error: "Forbidden"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResponseRecorder to record the response
			recorder := httptest.NewRecorder()
			// Call the function
			utils.WriteErrorResponse(recorder, tt.errorMessage, tt.statusCode)

			// Check the Content-Type header
			if recorder.Header().Get("Content-Type") != tt.expectedContentType {
				t.Errorf("expected Content-Type %s, got %s", tt.expectedContentType, recorder.Header().Get("Content-Type"))
			}

			// Check the status code
			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, recorder.Code)
			}

			// Check the response body
			var responseBody models.Response400
			if err := json.NewDecoder(recorder.Body).Decode(&responseBody); err != nil {
				t.Errorf("failed to decode response body: %v", err)
			}

			if responseBody != tt.expectedResponseBody {
				t.Errorf("expected response body %v, got %v", tt.expectedResponseBody, responseBody)
			}
		})
	}
}

func TestTemperatureHelper(t *testing.T) {
	tests := []struct {
		name          string
		actualPayload models.TempPostPayload
		expectedResp  models.TempPostResponse
	}{
		{
			name: "Temperature above threshold",
			actualPayload: models.TempPostPayload{
				DeviceId:    1234,
				EpochMS:     1721964434,
				Temperature: 95.0,
			},
			expectedResp: models.TempPostResponse{
				Overtemp:      true,
				DeviceId:      1234,
				FormattedTime: time.Unix(1721964434, 0).Format("2006/01/02 15:04:05"),
			},
		},
		{
			name: "Temperature below threshold",
			actualPayload: models.TempPostPayload{
				DeviceId:    1234,
				EpochMS:     1721964434,
				Temperature: 85.0,
			},
			expectedResp: models.TempPostResponse{
				Overtemp: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp models.TempPostResponse
			utils.TemperatureHelper(&tt.actualPayload, &resp)

			if resp.Overtemp != tt.expectedResp.Overtemp {
				t.Errorf("expected Overtemp %v, got %v", tt.expectedResp.Overtemp, resp.Overtemp)
			}
			if resp.DeviceId != tt.expectedResp.DeviceId {
				t.Errorf("expected DeviceId %d, got %d", tt.expectedResp.DeviceId, resp.DeviceId)
			}
			if resp.FormattedTime != tt.expectedResp.FormattedTime {
				t.Errorf("expected FormattedTime %s, got %s", tt.expectedResp.FormattedTime, resp.FormattedTime)
			}
		})
	}
}
