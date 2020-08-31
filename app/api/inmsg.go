package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gosdk/utils/logger"

	"message-queue/app/schema"
	"message-queue/app/services"
	"message-queue/pkg/app"
)

type InMsg struct {
	service services.InService
}

func NewInMsg(service services.InService) *InMsg {
	return &InMsg{service: service}
}

// Get List In Messages godoc
// @Tags In Messages
// @Summary get list in messages
// @Description get list in messages
// @Accept  json
// @Produce json
// @Param Query query schema.InMsgQueryParam true "Query"
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Header 200 {string} Token "qwerty"
// @Router /api/v1/in_messages [get]
func (o *InMsg) List(c *gin.Context) {
	var queryParam schema.InMsgQueryParam
	if err := c.Bind(&queryParam); err != nil {
		logger.Error("Failed to bind body, error: ", err)
		app.ResError(c, err, 400)
		return
	}

	rs, pageInfo, err := o.service.List(c, &queryParam)
	if err != nil {
		logger.Error("Failed to get list in messages, error: ", err)
		app.ResError(c, err, 400)
		return
	}

	res := schema.ResponsePaging{
		Data:   rs,
		Paging: pageInfo,
	}

	app.ResSuccess(c, res)
}
