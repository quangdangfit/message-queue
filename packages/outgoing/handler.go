package outgoing

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

type Handler interface {
	HandleMessage(message *OutMessage) (err error)
}

type handler struct {
	repo Repository
}

func NewHandler() Handler {
	return &handler{
		repo: NewRepository(),
	}
}

func (h *handler) HandleMessage(message *OutMessage) (
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
