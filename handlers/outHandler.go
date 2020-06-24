package handlers

import (
	"gomq/models"
	"gomq/repositories"

	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

type OutMessageHandler interface {
	HandleMessage(message *models.OutMessage) (err error)
}

type outHandler struct {
	repo repositories.OutMessageRepository
}

func NewOutMessageHandler() OutMessageHandler {
	return &outHandler{
		repo: repositories.NewOutMessageRepo(),
	}
}

func (s *outHandler) HandleMessage(message *models.OutMessage) (
	err error) {

	msg, err := s.repo.GetSingleOutMessage(bson.M{"id": message.ID})
	if msg != nil {
		err = s.repo.UpdateOutMessage(message)
		if err != nil {
			logger.Errorf("[Handle OutMsg] Failed to update msg %s, %s", message.ID, err)
			return err
		}

		logger.Infof("[Handle OutMsg] Updated msg %s", message.ID)
		return nil
	}

	err = s.repo.AddOutMessage(message)
	if err != nil {
		logger.Errorf("[Handle OutMsg] Failed to insert msg %s, %s", message.ID, err)
		return err
	}
	logger.Infof("[Handle OutMsg] Inserted msg %s", message.ID)
	return nil
}
