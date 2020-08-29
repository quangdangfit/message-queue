package api

import (
	"errors"

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
func (o *OutMsg) Publish(c *gin.Context) {
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

	message, err := o.prepareMessage(c, req)
	if err != nil {
		logger.Error("Failed to parse out message: ", err)
		app.ResError(c, err, 400)
		return
	}

	err = o.service.Publish(c, message)
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

// Update Out Message godoc
// @Tags Out Message
// @Summary api update out message
// @Description api update out message
// @Accept  json
// @Produce json
// @Param id path string true "Message ID"
// @Param Body body schema.OutMsgUpdateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Router /api/v1/out_messages/{id} [put]
func (o *OutMsg) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		err := errors.New("missing message id")
		logger.Error(err)
		app.ResError(c, err, 400)
		return
	}

	var bodyParam schema.OutMsgUpdateParam
	if err := c.Bind(&bodyParam); err != nil {
		logger.Error("Failed to bind body: ", err)
		app.ResError(c, err, 400)
		return
	}

	validate := validator.New()
	if err := validate.Validate(bodyParam); err != nil {
		logger.Error("Body is invalid: ", err)
		app.ResError(c, err, 400)
		return
	}

	rs, err := o.service.Update(c, id, &bodyParam)
	if err != nil {
		logger.Errorf("Failed to update out message %s, error: %s", id, err)
		app.ResError(c, err, 400)
	}

	app.ResSuccess(c, rs)
}

func (o *OutMsg) prepareMessage(c *gin.Context, body schema.OutMsgCreateParam) (
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
