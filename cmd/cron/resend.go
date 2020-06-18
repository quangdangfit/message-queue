package main

import (
	"gomq/dbs"
	"gomq/msgQueue"

	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

const ResendOutMessageLimit = 100

func main() {
	messages, _ := dbs.GetWaitOutMessage(ResendOutMessageLimit)
	if messages == nil {
		logger.Info("Not found any wait message!")
		return
	}

	pub := msgQueue.NewPublisher(false)
	for _, msg := range messages {
		err := pub.Publish(&msg, true)

		if err != nil {
			logger.Error("Failed to resend msg: ", msg.RoutingKey, msg.OriginModel, msg.OriginCode, err)
			msg.Logs = err.Error()
			msg.Status = dbs.OutMessageStatusFailed
		} else {
			msg.Status = dbs.OutMessageStatusSent
		}

		err = dbs.UpdateOutMessage(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, error: %s", msg.RoutingKey, msg.OriginModel, msg.OriginCode, err)
			continue
		}
	}
}
