package main

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"

	"gomq/packages/database"
	"gomq/packages/outgoing"
	"gomq/packages/queue"
)

const ResendOutMessageLimit = 100

func main() {
	repo := outgoing.NewRepository()

	query := bson.M{"status": database.OutMessageStatusWait}
	messages, _ := repo.GetOutMessages(query, ResendOutMessageLimit)
	if messages == nil {
		logger.Info("[Resend Message] Not found any wait message!")
		return
	}

	pub := queue.NewPublisher()

	logger.Infof("[Resend Message] Found %d wait messages!", len(*messages))
	for _, msg := range *messages {
		pub.Publish(&msg, true)
	}
	logger.Info("[Resend Message] Finish!")
}
