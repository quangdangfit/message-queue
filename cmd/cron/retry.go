package main

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"

	"gomq/app/models"
)

const (
	MaxRetryTimes       = 3
	RetryInMessageLimit = 100
)

func main() {
	repo := incoming.NewRepository()

	query := bson.M{"status": models.InMessageStatusWaitRetry}
	messages, _ := repo.GetInMessages(query, RetryInMessageLimit)
	if messages == nil {
		logger.Info("[Retry Message] Not found any wait_retry message!")
		return
	}

	inHandler := incoming.NewHandler()

	logger.Infof("[Retry Message] Found %d wait_retry messages!", len(*messages))
	for _, msg := range *messages {
		err := inHandler.HandleMessage(&msg, msg.RoutingKey.Name)
		if err == nil {
			continue
		}

		msg.Attempts += 1
		if msg.Attempts >= MaxRetryTimes {
			msg.Status = models.InMessageStatusFailed
		}
		err = repo.UpdateInMessage(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
		}
	}
	logger.Info("[Retry Message] Finish!")
}
