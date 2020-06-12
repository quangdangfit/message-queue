package service

import (
	"github.com/jinzhu/copier"
	"gomq/dbs"
	"gomq/models"
	"gomq/msgQueue"
	"net/http"
	"time"

	"transport/lib/utils/logger"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type Sender interface {
	PublishMessage(c echo.Context) (err error)
	parseMessage(msgRequest models.MessageRequest) (*models.OutMessage, error)
}

type sender struct {
	pub msgQueue.Publisher
}

func NewSender() Sender {
	pub := msgQueue.NewPublisher()
	return &sender{pub: pub}
}

func (s *sender) PublishMessage(c echo.Context) (err error) {
	var req models.MessageRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	validate := validator.New()
	if err = validate.Struct(req); err != nil {
		logger.Error("Publish: Bad request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	message, err := s.parseMessage(req)
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

func (s *sender) parseMessage(msgRequest models.MessageRequest) (
	*models.OutMessage, error) {
	message := models.OutMessage{}
	err := copier.Copy(&message, &msgRequest)

	if err != nil {
		return &message, err
	}

	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.Status = dbs.OutMessageStatusWait

	return &message, nil
}
