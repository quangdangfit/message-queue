package api

import (
	"net/http"

	"gomq/app/models"
	"gomq/app/queue"
	"gomq/utils"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/validator"
)

type Sender struct {
	pub queue.Publisher
}

func NewSender(pub queue.Publisher) *Sender {
	return &Sender{pub: pub}
}

// PublishMessage godoc
// @Summary publish message to amqp
// @Produce json
// @Body schema.OutMessageBodyParam
// @Security ApiKeyAuth
// @Success 200 {object} schema.OutMessageBodyParam
// @Router /api/v1/queue/messages [post]
func (s *Sender) PublishMessage(c echo.Context) (err error) {
	var req utils.MessageRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
	}

	validator := validator.New()
	if err = validator.Validate(req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
	}

	message, err := s.parseMessage(c, req)
	if err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
	}

	err = s.pub.Publish(message, true)
	if err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, utils.MsgResponse(utils.StatusBadRequest, nil))
	}

	return c.JSON(http.StatusOK, utils.MsgResponse(utils.StatusOK, nil))
}

func (s *Sender) parseMessage(c echo.Context, msgRequest utils.MessageRequest) (
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

func (s *Sender) getAPIKey(c echo.Context) string {
	return c.Request().Header.Get("X-Api-Key")
}
