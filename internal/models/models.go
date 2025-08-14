package models

import (
	"time"
	"github.com/google/uuid"
)

type APIResponse struct {
    ID           string      `json:"id"`
    Version      string      `json:"ver"`
    Timestamp    time.Time   `json:"ts"`
    Params       Params      `json:"params"`
    ResponseCode string      `json:"responseCode"`
    Result       interface{} `json:"result"`
}

type Params struct {
    MsgID     string `json:"msgid"`
}

type ApiRequest struct {
	ID           string      `json:"id"`
    Version      string      `json:"ver"`
    Timestamp    time.Time   `json:"ts"`
	Params       Params      `json:"params"`
	Request interface{} `json:request`
}

type Employees struct {
    ID        int     `gorm:"primaryKey;autoIncrement" json:"id"`
    FirstName string  `gorm:"size:50;not null" json:"first_name"`
    LastName  string  `gorm:"size:50;not null" json:"last_name"`
    Email     string  `gorm:"size:100;not null;unique" json:"email"`
    Salary    float64 `gorm:"type:numeric(10,2)" json:"salary"`
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