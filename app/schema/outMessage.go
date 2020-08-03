package schema

import "net/http"

type OutMessageQueryParam struct {
	RoutingKey  string `json:"routing_key,omitempty"`
	OriginCode  string `json:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty"`
	Status      string `json:"status,omitempty"`
}

type OutMessageBodyParam struct {
	RoutingKey  string      `json:"routing_key,omitempty" validate:"required"`
	Payload     interface{} `json:"payload,omitempty" validate:"required"`
	OriginCode  string      `json:"origin_code,omitempty"`
	OriginModel string      `json:"origin_model,omitempty"`
	Headers     http.Header `json:"headers,omitempty"`
}
