package service

import (
	"gomq/models"
	"gomq/msgHandler"
	"net/http"

	"transport/lib/utils/logger"

	"github.com/labstack/echo"
)

func PublishMessage(c echo.Context) (err error) {
	var req models.MessageRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("SCHEDULE: Bad request: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	validate := validator.New()
	if err = validate.Struct(req); err != nil {
		logger.Error("SCHEDULE: Bad request: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	sender := msgHandler.NewSender()
	if err = sender.Send(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{})
	}

	return c.JSON(http.StatusOK, nil)
}
