package msgHandler

import (
	"gomq/dbs"
	"gomq/models"
)

type OutMessageHandler interface {
	HandleMessage(message *models.OutMessage) (*models.OutMessage, error)
}

type outHandler struct{}

func NewOutMessageHandler() OutMessageHandler {
	return &outHandler{}
}

func (s *outHandler) HandleMessage(message *models.OutMessage) (
	*models.OutMessage, error) {
	dbs.AddOutMessage(message)
	return message, nil
}
