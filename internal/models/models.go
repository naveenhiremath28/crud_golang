package models

import (
	"time"
	"github.com/google/uuid"
)

type APIResponse struct {
    ID           string      `json:"id"`
    Version      string      `json:"ver"`
    Timestamp    time.Time      `json:"ts"`
    Params       Params      `json:"params"`
    ResponseCode string      `json:"responseCode"`
    Result       interface{} `json:"result"`
}

type Params struct {
    MsgID     string `json:"msgid"`
}

type ApiRequest struct {
	Request interface{} `json:request`
}

func GetApiResponse(id string, response_code string, result any) APIResponse {
	return APIResponse{
		ID: id,
		Version: "v1",
		Timestamp: time.Now(),
		Params: Params{
			MsgID: uuid.New().String(),
		},
		ResponseCode: response_code,
		Result: result,
	}
}