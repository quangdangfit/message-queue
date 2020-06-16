package models

import (
	"time"
)

type OutMessage struct {
	RoutingKey  string      `json:"routing_key,omitempty" bson:"routing_key,omitempty"`
	Payload     interface{} `json:"payload,omitempty" bson:"payload,omitempty"`
	OriginCode  string      `json:"origin_code,omitempty" bson:"origin_code,omitempty"`
	OriginModel string      `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
	Status      string      `json:"status,omitempty" bson:"status,omitempty"`
	Logs        string      `json:"logs,omitempty" bson:"logs,omitempty"`

	CreatedTime time.Time `json:"created_time" bson:"created_time"`
	UpdatedTime time.Time `json:"updated_time" bson:"updated_time"`
}

type InMessage struct {
	RoutingKey  RoutingKey  `json:"routing_key,omitempty" bson:"routing_key,omitempty"`
	Payload     interface{} `json:"payload,omitempty" bson:"payload,omitempty"`
	OriginCode  string      `json:"origin_code,omitempty" bson:"origin_code,omitempty"`
	OriginModel string      `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
	Status      string      `json:"status,omitempty" bson:"status,omitempty"`
	Logs        interface{} `json:"logs,omitempty" bson:"logs,omitempty"`
	Attempts    uint        `json:"attempts" bson:"attempts"`

	CreatedTime time.Time `json:"created_time" bson:"created_time"`
	UpdatedTime time.Time `json:"updated_time" bson:"updated_time"`
}
