package service

import (
	"github.com/jinzhu/copier"
	"gomq/dbs"
	"gomq/models"
	"gomq/msgQueue"
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"gopkg.in/go-playground/validator.v9"
)

type Sender interface {
	PublishMessage(c echo.Context) (err error)
	parseMessage(c echo.Context, msgRequest models.MessageRequest) (
		*models.OutMessage, error)
	getAPIKey(c echo.Context) string
}

type sender struct {
	pub msgQueue.Publisher
}

func NewSender() Sender {
	pub := msgQueue.NewPublisher(true)
	return &sender{pub: pub}
}

func (s *sender) PublishMessage(c echo.Context) (err error) {
	var req models.MessageRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	validator := validator.New()
	if err = validator.Struct(req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	message, err := s.parseMessage(c, req)
	if err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	err = s.pub.Publish(message, true)
	if err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *sender) parseMessage(c echo.Context, msgRequest models.MessageRequest) (
	*models.OutMessage, error) {
	message := models.OutMessage{}
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
