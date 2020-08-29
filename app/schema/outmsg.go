package schema

type OutMsgQueryParam struct {
	RoutingKey  string `json:"routing_key,omitempty" form:"routing_key,omitempty"`
	OriginCode  string `json:"origin_code,omitempty" form:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty" form:"origin_model,omitempty"`
	Status      string `json:"status,omitempty" form:"status,omitempty"`
	Page        int    `json:"-" form:"page,omitempty"`
	Limit       int    `json:"-" form:"limit,omitempty"`
}

type OutMsgCreateParam struct {
	RoutingKey  string      `json:"routing_key,omitempty" validate:"required" example:"routing.key"`
	Payload     interface{} `json:"payload,omitempty" validate:"required"`
	OriginCode  string      `json:"origin_code,omitempty" example:"code"`
	OriginModel string      `json:"origin_model,omitempty" example:"model"`
}

type OutMsgUpdateParam struct {
	Status string `json:"status,omitempty" validate:"omitempty,oneof=wait canceled sent"`
}
