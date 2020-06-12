package utils

import (
	"encoding/json"
	"net/http"
)

type objResponse struct {
	HttpCode        int    `json:"http_code"`
	Code            string `json:"code"`
	Message         string `json:"message"`
	OriginalMessage string `json:"original_message,omitempty"`
}

func ParseError(res http.Response) *objResponse {
	error := objResponse{}
	json.NewDecoder(res.Body).Decode(&error)
	return &error
}
