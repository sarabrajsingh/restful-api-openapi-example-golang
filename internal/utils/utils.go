package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sarabrajsingh/restful-openapi/internal/models"
)

// a helper function that marshalls an error string into a proper response object to be
// written to the http response writer
func WriteErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	errorResponse := models.Response400{
		Error: errorMessage,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

type BodyReaderFunc func(io.Reader) ([]byte, error)

func DefaultBodyReader(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func TemperatureHelper(actual *models.TempPostPayload, response *models.TempPostResponse) {
	// parse out the meat and potatoes
	if actual.Temperature >= 90.00 {
		response.Overtemp = true
		response.DeviceId = actual.DeviceId
		// convert epoch into an epochInt and then into a time object
		epochInt := int64(actual.EpochMS)
		t := time.Unix(epochInt, 0)
		formattedTime := t.Format("2006/01/02 15:04:05")
		response.FormattedTime = formattedTime
	} else {
		response.Overtemp = false
	}
}

func PayloadParserHelper(payloadData string) (*models.TempPostPayload, error) {
	data := strings.Split(payloadData, ":")

	if len(data) != 4 {
		return nil, fmt.Errorf("invalid number of arguments in request body")
	}

	deviceId, err := strconv.ParseInt(data[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse device_id=%s to an int32", data[0])
	}

	epochMS, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse epochMS=%s to an int64", data[1])
	}

	if string(data[2]) != string("'Temperature'") {
		return nil, fmt.Errorf("temperature key is mislabelled: %s", data[2])
	}

	temperature, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse temperature=%s to a float64", data[3])
	}

	// these casts may be unncessary
	return &models.TempPostPayload{
		DeviceId:    int32(deviceId),
		EpochMS:     int64(epochMS),
		Temperature: float64(temperature),
	}, nil
}
