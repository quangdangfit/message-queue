package main

import (
	"gomq/dbs"
	"gomq/packages/incomming"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

const (
	MaxRetryTimes       = 3
	RetryInMessageLimit = 100
)

func main() {
	repo := incomming.NewInMessageRepo()

	query := bson.M{"status": dbs.InMessageStatusWaitRetry}
	messages, _ := repo.GetInMessages(query, RetryInMessageLimit)
	if messages == nil {
		logger.Info("[Retry Message] Not found any wait_retry message!")
		return
	}

	inHandler := incomming.NewInMessageHandler()

	logger.Infof("[Retry Message] Found %d wait_retry messages!", len(*messages))
	for _, msg := range *messages {
		err := inHandler.HandleMessage(&msg, msg.RoutingKey.Name)
		if err == nil {
			continue
		}

		msg.Attempts += 1
		if msg.Attempts >= MaxRetryTimes {
			msg.Status = dbs.InMessageStatusFailed
		}
		err = repo.UpdateInMessage(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
		}
	}
	logger.Info("[Retry Message] Finish!")
}
