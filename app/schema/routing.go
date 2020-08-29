package schema

type RoutingKeyQueryParam struct {
	Group string `json:"group,omitempty" form:"group,omitempty"`
	Name  string `json:"name,omitempty" form:"name,omitempty"`
	Value uint   `json:"value,omitempty" form:"value,omitempty"`
	Page  int    `json:"-" form:"page,omitempty"`
	Limit int    `json:"-" form:"limit,omitempty"`
}
