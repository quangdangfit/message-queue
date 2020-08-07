package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gomq/app/models"
	"gomq/app/services"
	"gomq/utils"

	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/validator"
)

type Sender struct {
	service services.OutService
}

func NewSender(service services.OutService) *Sender {
	return &Sender{service: service}
}

// PublishMessage godoc
// @Tags Out Messages
// @Summary publish message to amqp
// @Description api publish out message to amqp
// @Accept  json
// @Produce json
// @Param Body body schema.OutMessageBodyParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} schema.OutMessageBodyParam
// @Header 200 {string} Token "qwerty"
// @Router /api/v1/queue/messages [post]
func (s *Sender) PublishMessage(c *gin.Context) {
	var req utils.MessageRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
		return
	}

	validate := validator.New()
	if err := validate.Validate(req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
		return
	}

	message, err := s.parseMessage(c, req)
	if err != nil {
		logger.Error("Publish: Bad request: ", err)
		c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
		return
	}

	err = s.service.Publish(c, message)
	if err != nil {
		logger.Error("Publish: Bad request: ", err)
		c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
		return
	}

	c.JSON(http.StatusOK, utils.MsgResponse(utils.StatusOK, nil))
}

func (s *Sender) parseMessage(c *gin.Context, msgRequest utils.MessageRequest) (
	*models.OutMessage, error) {
	message := models.OutMessage{}
	err := copier.Copy(&message, &msgRequest)

	if err != nil {
		return &message, err
	}
	message.Status = models.OutMessageStatusWait
	message.APIKey = s.getAPIKey(c)

	return &message, nil
}

func (s *Sender) getAPIKey(c *gin.Context) string {
	return c.Request.Header.Get("X-Api-Key")
}
