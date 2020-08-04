package schema

type OutMessageQueryParam struct {
	RoutingKey  string `json:"routing_key,omitempty"`
	OriginCode  string `json:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty"`
	Status      string `json:"status,omitempty"`
}

type OutMessageBodyParam struct {
	RoutingKey  string      `json:"routing_key,omitempty" validate:"required" example:"routing.key"`
	Payload     interface{} `json:"payload,omitempty" validate:"required"`
	OriginCode  string      `json:"origin_code,omitempty" example:"code"`
	OriginModel string      `json:"origin_model,omitempty" example:"model"`
}
