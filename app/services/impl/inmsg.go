package impl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/spf13/viper"

	"message-queue/app/models"
	"message-queue/app/queue"
	"message-queue/app/repositories"
	"message-queue/app/schema"
	"message-queue/app/services"
	"message-queue/pkg/utils"
)

const (
	RequestTimeout         = 60
	DefaultMaxRetryTimes   = 3
	DefaultConsumerThreads = 10
)

type inService struct {
	inRepo      repositories.InRepository
	routingRepo repositories.RoutingRepository

	consumer        queue.Consumer
	consumerThreads int
}

func NewInService(inRepo repositories.InRepository, routingRepo repositories.RoutingRepository,
	consumer queue.Consumer) services.InService {

	r := inService{
		inRepo:          inRepo,
		routingRepo:     routingRepo,
		consumer:        consumer,
		consumerThreads: DefaultConsumerThreads,
	}
	return &r
}

func (i *inService) Consume() {
	msgChan := i.consumer.Consume()
	logger.Infof("Run %d threads to consume messages", i.consumerThreads)
	for index := 0; index <= i.consumerThreads; index++ {
		for msg := range msgChan {
			i.handle(msg, msg.RoutingKey.Name)
			i.inRepo.Create(msg)
		}
	}
}

func (i *inService) CronRetry(limit int) error {
	query := schema.InMsgQueryParam{
		Status: models.InMessageStatusWaitRetry,
	}

	messages, _ := i.inRepo.List(&query, limit)
	if messages == nil {
		logger.Info("[Retry Message] Not found any wait_retry message!")
		return nil
	}

	logger.Infof("[Retry Message] Found %d wait_retry messages!", len(*messages))
	for _, msg := range *messages {
		err := i.handle(&msg, msg.RoutingKey.Name)
		if err == nil {
			continue
		}

		msg.Attempts += 1
		if msg.Attempts >= i.getMaxRetryTimes() {
			msg.Status = models.InMessageStatusFailed
		}
		err = i.inRepo.Update(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, error: %s",
				msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
		}
	}
	logger.Info("[Retry Message] Finish!")

	return nil
}

func (i *inService) CronRetryPrevious(limit int) error {
	query := schema.InMsgQueryParam{
		Status: models.InMessageStatusWaitPrevMsg,
	}
	messages, _ := i.inRepo.List(&query, limit)
	if messages == nil {
		logger.Info("[Retry Prev Message] Not found any wait_prev message!")
		return nil
	}

	logger.Infof("[Retry Prev Message] Found %d wait_prev messages!", len(*messages))
	for _, msg := range *messages {
		query := schema.InMsgQueryParam{
			RoutingGroup: msg.RoutingKey.Group,
			RoutingValue: msg.RoutingKey.Value - 1,
		}
		prevMsg, err := i.inRepo.Get(&query)
		if (prevMsg == nil && msg.RoutingKey.Value != 1) ||
			(prevMsg != nil && prevMsg.Status != models.InMessageStatusSuccess &&
				prevMsg.Status != models.InMessageStatusCanceled) {

			logger.Infof("[Retry Prev Message] Ignore message %s!", msg.ID)
			continue
		}

		err = i.handle(&msg, msg.RoutingKey.Name)
		if err == nil {
			continue
		}

		msg.Attempts += 1
		if msg.Attempts >= i.getMaxRetryTimes() {
			msg.Status = models.InMessageStatusFailed
		}
		err = i.inRepo.Update(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, "+
				"error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
		}
	}
	logger.Info("[Retry Prev Message] Finish!")

	return nil
}

func (i *inService) handle(message *models.InMessage, routingKey string) error {
	query := schema.RoutingQueryParam{
		Name: routingKey,
	}
	inRoutingKey, err := i.routingRepo.Get(&query)
	if err != nil {
		message.Status = models.InMessageStatusInvalid
		message.Logs = append(message.Logs, utils.ParseLogs(err))
		logger.Error("Cannot find routing key ", err)
		return err
	}
	message.RoutingKey = *inRoutingKey

	prevMsg, _ := i.getPrevMessage(message)
	if (prevMsg == nil && message.RoutingKey.Value != 1) ||
		(prevMsg != nil && prevMsg.Status != models.InMessageStatusSuccess && prevMsg.Status != models.InMessageStatusCanceled) {
		message.Status = models.InMessageStatusWaitPrevMsg
		logger.Warn("Set message to WAIT_PREV_MESSAGE")
		return nil
	}

	res, err := i.callAPI(message)
	if err != nil {
		message.Status = models.InMessageStatusWaitRetry
		message.Logs = append(message.Logs, utils.ParseLogs(err))
		return err
	}

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusUnauthorized {
		message.Status = models.InMessageStatusWaitRetry
		err = errors.New(fmt.Sprintf("failed to call API %s", res.Status))
		message.Logs = append(message.Logs, utils.ParseLogs(res))
		return err
	}

	if res.StatusCode != http.StatusOK {
		message.Status = models.InMessageStatusWaitRetry
		err = errors.New("failed to call API")
		message.Logs = append(message.Logs, utils.ParseLogs(res))
		return err
	}

	message.Status = models.InMessageStatusSuccess
	message.Logs = append(message.Logs, utils.ParseLogs(res))

	return nil
}

func (i *inService) storeMessage(message *models.InMessage) (err error) {
	return i.inRepo.Upsert(message)
}

func (i *inService) callAPI(message *models.InMessage) (*http.Response, error) {
	routingKey := message.RoutingKey

	bytesPayload, _ := json.Marshal(message.Payload)
	req, _ := http.NewRequest(
		routingKey.APIMethod, routingKey.APIUrl, bytes.NewBuffer(bytesPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "ahsfishdi"))
	req.Header.Set("x-api-key", message.APIKey)

	client := http.Client{
		Timeout: RequestTimeout * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		logger.Errorf("Failed to send request to %s, %s", routingKey.APIUrl, err)
		return res, err
	}

	return res, nil
}

func (i *inService) getPrevMessage(message *models.InMessage) (*models.InMessage, error) {
	// Get previous routing
	routingQuery := schema.RoutingQueryParam{
		Group: message.RoutingKey.Group,
		Value: message.RoutingKey.Value - 1,
	}
	prevRouting, err := i.routingRepo.Get(&routingQuery)
	if err != nil || prevRouting == nil {
		return nil, err
	}

	query := schema.InMsgQueryParam{
		OriginModel: message.OriginModel,
		OriginCode:  message.OriginCode,
		RoutingKey:  prevRouting.Name,
	}
	return i.inRepo.Get(&query)
}

func (i *inService) getMaxRetryTimes() uint {
	retryTimes := viper.GetUint("ts_rabbit.max_retry_times")

	if retryTimes <= 0 {
		retryTimes = DefaultMaxRetryTimes
	}

	return retryTimes
}
