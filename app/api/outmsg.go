package api

import (
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/validator"

	"message-queue/app/models"
	"message-queue/app/schema"
	"message-queue/app/services"
	"message-queue/pkg/app"
)

type OutMsg struct {
	service services.OutService
}

func NewSender(service services.OutService) *OutMsg {
	return &OutMsg{service: service}
}

// Publish Message godoc
// @Tags Out Messages
// @Summary publish message to amqp
// @Description api publish out message to amqp
// @Accept  json
// @Produce json
// @Param Body body schema.OutMsgCreateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Header 200 {string} Token "qwerty"
// @Router /api/v1/out_messages [post]
func (s *OutMsg) Publish(c *gin.Context) {
	var req schema.OutMsgCreateParam
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind body: ", err)
		app.ResError(c, err, 400)
		return
	}

	validate := validator.New()
	if err := validate.Validate(req); err != nil {
		logger.Error("Body is invalid: ", err)
		app.ResError(c, err, 400)
		return
	}

	message, err := s.prepareMessage(c, req)
	if err != nil {
		logger.Error("Failed to parse out message: ", err)
		app.ResError(c, err, 400)
		return
	}

	err = s.service.Publish(c, message)
	if err != nil {
		logger.Error("Failed to publish message: ", err)
		app.ResError(c, err, 400)
		return
	}

	app.ResOK(c)
}

// Get List Out Messages godoc
// @Tags Out Messages
// @Summary get list out messages
// @Description get list out messages
// @Accept  json
// @Produce json
// @Param Query query schema.OutMsgQueryParam true "Query"
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Header 200 {string} Token "qwerty"
// @Router /api/v1/out_messages [get]
func (o *OutMsg) List(c *gin.Context) {
	var queryParam schema.OutMsgQueryParam
	if err := c.Bind(&queryParam); err != nil {
		logger.Error("Failed to bind body, error: ", err)
		app.ResError(c, err, 400)
		return
	}

	rs, pageInfo, err := o.service.List(c, &queryParam)
	if err != nil {
		logger.Error("Failed to get list out messages, error: ", err)
		app.ResError(c, err, 400)
		return
	}

	res := schema.ResponsePaging{
		Data:   rs,
		Paging: pageInfo,
	}

	app.ResSuccess(c, res)
}

func (s *OutMsg) prepareMessage(c *gin.Context, body schema.OutMsgCreateParam) (
	*models.OutMessage, error) {
	message := models.OutMessage{}
	err := copier.Copy(&message, &body)

	if err != nil {
		return &message, err
	}
	message.Status = models.OutMessageStatusWait
	message.APIKey = c.Request.Header.Get("X-Api-Key")

	return &message, nil
}
