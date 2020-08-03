package services

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"

	"gomq/app/models"
	"gomq/app/repositories"
)

type handler struct {
	repo repositories.OutMessageRepository
}

func NewOutMessageService(repo repositories.OutMessageRepository) OutMessageService {
	return &handler{
		repo: repo,
	}
}

func (h *handler) HandleMessage(message *models.OutMessage) (
	err error) {

	msg, err := h.repo.GetSingleOutMessage(bson.M{"id": message.ID})
	if msg != nil {
		err = h.repo.UpdateOutMessage(message)
		if err != nil {
			logger.Errorf("[Handle OutMsg] Failed to update msg %s, %s", message.ID, err)
			return err
		}

		logger.Infof("[Handle OutMsg] Updated msg %s", message.ID)
		return nil
	}

	err = h.repo.AddOutMessage(message)
	if err != nil {
		logger.Errorf("[Handle OutMsg] Failed to insert msg %s, %s", message.ID, err)
		return err
	}
	logger.Infof("[Handle OutMsg] Inserted msg %s", message.ID)
	return nil
}
