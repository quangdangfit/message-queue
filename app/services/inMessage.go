package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"

	dbs "gomq/app/database"
	"gomq/app/models"
	"gomq/app/repositories"
	"gomq/utils"
)

const (
	RequestTimeout = 60
)

type inService struct {
	inMsgRepo   repositories.InMessageRepository
	routingRepo repositories.RoutingRepository
}

func NewInMessageService(inMsgRepo repositories.InMessageRepository, routingRepo repositories.RoutingRepository) InMessageService {
	r := inService{
		inMsgRepo:   inMsgRepo,
		routingRepo: routingRepo,
	}
	return &r
}

func (i *inService) HandleMessage(message *models.InMessage, routingKey string) error {

	defer i.storeMessage(message)

	query := bson.M{"name": routingKey}
	inRoutingKey, err := i.routingRepo.GetRoutingKey(query)
	if err != nil {
		message.Status = dbs.InMessageStatusInvalid
		message.Logs = append(message.Logs, utils.ParseLog(err))
		logger.Error("Cannot find routing key ", err)
		return err
	}
	message.RoutingKey = *inRoutingKey

	prevRoutingKey, _ := i.routingRepo.GetPreviousRoutingKey(message.RoutingKey)
	if prevRoutingKey != nil {
		prevMsg, _ := i.getPreviousMessage(*message, prevRoutingKey.Name)

		if prevMsg == nil || (prevMsg.Status != dbs.InMessageStatusSuccess &&
			prevMsg.Status != dbs.InMessageStatusCanceled) {

			message.Status = dbs.InMessageStatusWaitPrevMsg

			logger.Warn("Set message to WAIT_PREV_MESSAGE")
			return nil
		}
	}

	res, err := i.callAPI(message)
	if err != nil {
		message.Status = dbs.InMessageStatusWaitRetry
		message.Logs = append(message.Logs, utils.ParseLog(err))
		return err
	}

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusUnauthorized {
		message.Status = dbs.InMessageStatusWaitRetry
		err = errors.New(fmt.Sprintf("failed to call API %s", res.Status))
		message.Logs = append(message.Logs, utils.ParseLog(res))
		return err
	}

	if res.StatusCode != http.StatusOK {
		message.Status = dbs.InMessageStatusWaitRetry
		err = errors.New("failed to call API")
		message.Logs = append(message.Logs, utils.ParseLog(res))
		return err
	}

	message.Status = dbs.InMessageStatusSuccess
	message.Logs = append(message.Logs, utils.ParseLog(res))

	return nil
}

func (i *inService) storeMessage(message *models.InMessage) (err error) {
	return i.inMsgRepo.UpsertInMessage(message)
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

func (i *inService) getPreviousMessage(message models.InMessage, routingKey string) (
	*models.InMessage, error) {

	query := bson.M{
		"origin_model":     message.OriginModel,
		"origin_code":      message.OriginCode,
		"routing_key.name": routingKey,
	}
	return i.inMsgRepo.GetSingleInMessage(query)
}
