package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type APIResponse struct {
	ID           string      `json:"id"`
	Version      string      `json:"ver"`
	Timestamp    string      `json:"ts"`
	Params       Params      `json:"params"`
	ResponseCode string      `json:"responseCode"`
	Result       interface{} `json:"result"`
}

type Params struct {
	MsgID string `json:"msgid"`
}

type ApiRequest struct {
	ID        string          `json:"id"`
	Version   string          `json:"ver"`
	Timestamp string          `json:"ts"`
	Params    Params          `json:"params"`
	Request   json.RawMessage `json:"request"`
}
	
type Employees struct {
	ID            string  `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName     string  `gorm:"size:50;not null" json:"first_name"`
	LastName      string  `gorm:"size:50;not null" json:"last_name"`
	Email         string  `gorm:"size:100;not null;unique" json:"email"`
	Salary        float64 `gorm:"type:numeric(10,2)" json:"salary"`
	Mobile        string  `gorm:"size:15;not null" json:"mobile"`
	VaultEntityID string  `gorm:"size:100" json:"vault_entity_id"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest[T any] struct {
	ID        string `json:"id"`
	Version   string `json:"ver"`
	Timestamp string `json:"ts"`
	Params    Params `json:"params"`
	Request   T      `json:"request"`
}

func GetApiResponse(id string, response_code string, result any) APIResponse {
	return APIResponse{
		ID:        id,
		Version:   "v1",
		Timestamp: time.Now().Format(time.RFC3339),
		Params: Params{
			MsgID: uuid.New().String(),
		},
		ResponseCode: response_code,
		Result:       result,
	}
}
