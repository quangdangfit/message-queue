package utils

import (
	"encoding/json"
	"net/http"
)

type objLogs struct {
	Status     string      `json:"status,omitempty" bson:"status,omitempty"`
	StatusCode int         `json:"status_code,omitempty" bson:"status_code,omitempty"`
	Body       interface{} `json:"body,omitempty" bson:"body,omitempty"`
	Error      string      `json:"error,omitempty" bson:"error,omitempty"`
}

func ParseLog(err interface{}) *objLogs {
	logObject := objLogs{}
	switch v := err.(type) {
	case *http.Response:
		var body interface{}
		json.NewDecoder(v.Body).Decode(&body)

		logObject.Status = v.Status
		logObject.StatusCode = v.StatusCode
		logObject.Body = body
		//break
	case error:
		logObject.Error = v.Error()
	case string:
		logObject.Error = v
	}

	return &logObject
}
