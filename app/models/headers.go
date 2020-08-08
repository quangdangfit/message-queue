package models

type Headers struct {
	OriginCode  string `json:"origin_code,omitempty" bson:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
	APIKey      string `json:"api_key,omitempty" bson:"api_key,omitempty"`
}
