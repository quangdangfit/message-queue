package schema

type InMsgQueryParam struct {
	RoutingKey   string `json:"routing_key.name,omitempty" form:"routing_key.name,omitempty"`
	RoutingGroup string `json:"routing_key.group,omitempty" form:"routing_key.group,omitempty"`
	RoutingValue uint   `json:"routing_key.value,omitempty" form:"routing_key.value,omitempty"`
	OriginCode   string `json:"origin_code,omitempty" form:"origin_code,omitempty"`
	OriginModel  string `json:"origin_model,omitempty" form:"origin_model,omitempty"`
	Status       string `json:"status,omitempty" form:"status,omitempty"`
	Page         int    `json:"-" form:"page,omitempty"`
	Limit        int    `json:"-" form:"limit,omitempty"`
}
