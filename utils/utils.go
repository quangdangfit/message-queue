package utils

import (
	"encoding/json"
	"net/http"
)

type objResponse struct {
	HttpCode int    `json:"http_code,omitempty" bson:"http_code,omitempty"`
	Code     string `json:"code,omitempty" bson:"code,omitempty"`
	Message  string `json:"message,omitempty" bson:"message,omitempty"`
}

func ParseError(err interface{}) *objResponse {
	errObj := objResponse{}
	switch v := err.(type) {
	case http.Response:
		json.NewDecoder(v.Body).Decode(&errObj)
		break
	case error:
		errObj.Message = v.Error()
	}

	return &errObj
}
