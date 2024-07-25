package models

type GetErrorsResponse struct {
	Errors []string `json:"errors"`
}

type Response400 struct {
	Error string `json:"error"`
}

type TempPostBody struct {
	Data string `json:"data"`
}

type TempPostPayload struct {
	DeviceId    int32
	EpochMS     int64
	Temperature float64
}

type TempPostResponse struct {
	Overtemp      bool   `json:"overtemp"`
	DeviceId      int32  `json:"device_id,omitempty"`
	FormattedTime string `json:"formatted_time,omitempty"`
}
