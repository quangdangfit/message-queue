package services

import (
	"github.com/quangdangfit/gosdk/utils/logger"

	"gomq/app/models"
	"gomq/app/repositories"
	"gomq/app/schema"
)

type outService struct {
	repo repositories.OutMessageRepository
}

func NewOutMessageService(repo repositories.OutMessageRepository) OutMessageService {
	return &outService{
		repo: repo,
	}
}

func (h *outService) HandleMessage(message *models.OutMessage) (err error) {
	msg, err := h.repo.GetOutMessageByID(message.ID)
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

func (o *outService) GetOutMessages(query *schema.OutMessageQueryParam, limit int) (*[]models.OutMessage, error) {
	msg, err := o.repo.GetOutMessages(query, limit)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
