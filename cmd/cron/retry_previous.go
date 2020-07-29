package main

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"

	dbs "gomq/packages/database"
	"gomq/packages/incoming"
)

const (
	RetryPrevInMessageLimit = 100
	MaxPrevRetryTimes       = 3
)

func main() {
	repo := incoming.NewRepository()

	query := bson.M{"status": dbs.InMessageStatusWaitPrevMsg}
	messages, _ := repo.GetInMessages(query, RetryPrevInMessageLimit)
	if messages == nil {
		logger.Info("[Retry Prev Message] Not found any wait_prev message!")
		return
	}

	inHandler := incoming.NewHandler()

	logger.Infof("[Retry Prev Message] Found %d wait_prev messages!", len(*messages))
	for _, msg := range *messages {
		query := bson.M{
			"routing_key.group": msg.RoutingKey.Group,
			"routing_key.value": msg.RoutingKey.Value - 1,
		}
		prevMsg, err := repo.GetSingleInMessage(query)
		if (prevMsg == nil && msg.RoutingKey.Value != 1) ||
			(prevMsg != nil && prevMsg.Status != dbs.InMessageStatusSuccess &&
				prevMsg.Status != dbs.InMessageStatusCanceled) {

			logger.Infof("[Retry Prev Message] Ignore message %s!", msg.ID)
			continue
		}

		err = inHandler.HandleMessage(&msg, msg.RoutingKey.Name)
		if err == nil {
			continue
		}

		msg.Attempts += 1
		if msg.Attempts >= MaxPrevRetryTimes {
			msg.Status = dbs.InMessageStatusFailed
		}
		err = repo.UpdateInMessage(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
		}
	}
	logger.Info("[Retry Prev Message] Finish!")
}
