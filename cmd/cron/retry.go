package main

import (
	"gomq/dbs"
	"gomq/handlers"
	"gomq/utils"

	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

const (
	MaxRetryTimes       = 3
	RetryInMessageLimit = 100
)

func main() {
	messages, _ := dbs.GetInMessageByStatus(dbs.InMessageStatusWaitRetry, RetryInMessageLimit)
	if messages == nil {
		logger.Info("Not found any wait retry message!")
		return
	}

	inHandler := handlers.NewInMessageHandler(false)
	for _, msg := range messages {
		_, err := inHandler.HandleMessage(&msg, msg.RoutingKey.Name)
		if err != nil {
			logger.Error("Failed to retry msg: ", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
			msg.Logs = append(msg.Logs, utils.ParseError(err))
			msg.Attempts += 1

			msg.Status = dbs.InMessageStatusWaitRetry
			if msg.Attempts > MaxRetryTimes {
				msg.Status = dbs.InMessageStatusFailed
			}
		} else {
			msg.Status = dbs.InMessageStatusSuccess
		}

		err = dbs.UpdateInMessage(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
			continue
		}
	}
}
