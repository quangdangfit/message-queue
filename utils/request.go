package utils

type MessageRequest struct {
	RoutingKey  string      `json:"routing_key,omitempty" bson:"routing_key,omitempty" validate:"required"`
	Payload     interface{} `json:"payload,omitempty" bson:"payload,omitempty" validate:"required"`
	OriginCode  string      `json:"origin_code,omitempty" bson:"origin_code,omitempty"`
	OriginModel string      `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
}
