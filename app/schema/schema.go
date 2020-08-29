package schema

import "github.com/quangdangfit/gosdk/utils/paging"

type ResponsePaging struct {
	Paging *paging.Paging `json:"paging"`
	Data   interface{}    `json:"data"`
}
