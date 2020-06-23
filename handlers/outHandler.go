package handlers

import (
	"gomq/dbs"
	"gomq/models"
)

type OutMessageHandler interface {
	HandleMessage(message *models.OutMessage, store bool) (*models.OutMessage, error)
}

type outHandler struct{}

func NewOutMessageHandler() OutMessageHandler {
	return &outHandler{}
}

func (s *outHandler) HandleMessage(message *models.OutMessage, store bool) (
	*models.OutMessage, error) {

	if store {
		dbs.AddOutMessage(message)
	}
	return message, nil
}
