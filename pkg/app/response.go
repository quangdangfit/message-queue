package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

var OK = "Success"

func ResError(c *gin.Context, err error, status int) {
	ResJSON(c, status, Response{Code: status, Msg: err.Error()})
}
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK,
		Response{
			Code: http.StatusOK,
			Msg:  OK,
			Data: v,
		})
}

func ResOK(c *gin.Context) {
	ResSuccess(c, nil)
}

func ResJSON(c *gin.Context, httpCode int, res Response) {
	c.JSON(httpCode, res)
	c.AbortWithStatus(httpCode)
	return
}
