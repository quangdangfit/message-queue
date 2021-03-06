package schema

type RoutingQueryParam struct {
	Group string `json:"group,omitempty" form:"group,omitempty"`
	Name  string `json:"name,omitempty" form:"name,omitempty"`
	Value uint   `json:"value,omitempty" form:"value,omitempty"`
	Page  int    `json:"-" form:"page,omitempty"`
	Limit int    `json:"-" form:"limit,omitempty"`
}

type RoutingCreateParam struct {
	Name      string `json:"name,omitempty" validate:"required"`
	Group     string `json:"group,omitempty" validate:"required"`
	Value     uint   `json:"value,omitempty" validate:"required,gt=0"`
	APIMethod string `json:"api_method,omitempty" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	APIUrl    string `json:"api_url,omitempty" validate:"required,url"`
}

type RoutingUpdateParam struct {
	Name      string `json:"name,omitempty"`
	Group     string `json:"group,omitempty"`
	Value     uint   `json:"value,omitempty" validate:"gt=0"`
	APIMethod string `json:"api_method,omitempty" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	APIUrl    string `json:"api_url,omitempty" validate:"omitempty,url"`
}
