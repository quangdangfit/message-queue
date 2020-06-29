package utils

import (
	"encoding/json"
	"net/http"
)

type objLogs struct {
	Status   string `json:"status,omitempty" bson:"status,omitempty"`
	HttpCode int    `json:"http_code,omitempty" bson:"http_code,omitempty"`
	Code     string `json:"code,omitempty" bson:"code,omitempty"`
	Message  string `json:"message,omitempty" bson:"message,omitempty"`
}

func ParseError(err interface{}) *objLogs {
	errObj := objLogs{}
	switch v := err.(type) {
	case http.Response:
		json.NewDecoder(v.Body).Decode(&errObj)
		break
	case error:
		errObj.Message = v.Error()
	case string:
		errObj.Message = v
	}

	return &errObj
}
