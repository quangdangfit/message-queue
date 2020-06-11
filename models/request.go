package models

type MessageRequest struct {
	RoutingKey  string      `json:"routing_key,omitempty" bson:"routing_key,omitempty" validate:"required"`
	Payload     interface{} `json:"payload,omitempty" bson:"payload,omitempty" validate:"required"`
	OriginID    string      `json:"origin_id,omitempty" bson:"origin_id,omitempty"`
	OriginModel string      `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
}
