package services

import (
	"net/http"

	dbs "gomq/packages/database"
	"gomq/packages/outgoing"
	"gomq/packages/queue"
	"gomq/utils"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/validator"
)

type Sender interface {
	PublishMessage(c echo.Context) (err error)
	parseMessage(c echo.Context, msgRequest utils.MessageRequest) (
		*outgoing.OutMessage, error)
	getAPIKey(c echo.Context) string
}

type sender struct {
	pub queue.Publisher
}

func NewSender() Sender {
	pub := queue.NewPublisher()
	return &sender{pub: pub}
}

func (s *sender) PublishMessage(c echo.Context) (err error) {
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

func (s *sender) parseMessage(c echo.Context, msgRequest utils.MessageRequest) (
	*outgoing.OutMessage, error) {
	message := outgoing.OutMessage{}
	err := copier.Copy(&message, &msgRequest)

	if err != nil {
		return &message, err
	}
	message.Status = dbs.OutMessageStatusWait
	message.APIKey = s.getAPIKey(c)

	return &message, nil
}

func (s *sender) getAPIKey(c echo.Context) string {
	return c.Request().Header.Get("X-Api-Key")
}
