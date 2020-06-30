package incomming

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gomq/dbs"
	"gomq/packages/inrouting"
	"gomq/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"

	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

const (
	RequestTimeout = 60
)

type InMessageHandler interface {
	HandleMessage(message *InMessage, routingKey string) error
	storeMessage(message *InMessage) (err error)
	callAPI(message *InMessage) (*http.Response, error)
}

type inHandler struct {
	repo InMessageRepository
}

func NewInMessageHandler() InMessageHandler {
	r := inHandler{
		repo: NewInMessageRepo(),
	}
	return &r
}

func (r *inHandler) HandleMessage(message *InMessage, routingKey string) error {

	defer r.storeMessage(message)

	routingRepo := inrouting.NewRoutingKeyRepo()
	inRoutingKey, err := routingRepo.GetRoutingKey(routingKey)
	if err != nil {
		message.Status = dbs.InMessageStatusInvalid
		message.Logs = append(message.Logs, utils.ParseError(err))
		logger.Error("Cannot find routing key ", err)
		return err
	}
	message.RoutingKey = *inRoutingKey

	prevRoutingKey, _ := routingRepo.GetPreviousRoutingKey(message.RoutingKey)
	if prevRoutingKey != nil {
		prevMsg, _ := r.getPreviousMessage(*message, prevRoutingKey.Name)

		if prevMsg == nil || (prevMsg.Status != dbs.InMessageStatusSuccess &&
			prevMsg.Status != dbs.InMessageStatusCanceled) {

			message.Status = dbs.InMessageStatusWaitPrevMsg

			logger.Warn("Set message to WAIT_PREV_MESSAGE")
			return nil
		}
	}

	res, err := r.callAPI(message)
	if err != nil {
		message.Status = dbs.InMessageStatusWaitRetry
		message.Logs = append(message.Logs, utils.ParseError(err))
		return err
	}

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusUnauthorized {
		message.Status = dbs.InMessageStatusWaitRetry
		err = errors.New(fmt.Sprintf("failed to call API %s", res.Status))
		message.Logs = append(message.Logs, utils.ParseError(res.Status))
		return err
	}

	if res.StatusCode != http.StatusOK {
		message.Status = dbs.InMessageStatusWaitRetry
		err = errors.New("failed to call API")
		message.Logs = append(message.Logs, utils.ParseError(res))
		return err
	}

	message.Status = dbs.InMessageStatusSuccess
	message.Logs = append(message.Logs, utils.ParseError(res))

	return nil
}

func (r *inHandler) storeMessage(message *InMessage) (err error) {
	msg, err := r.repo.GetSingleInMessage(bson.M{"id": message.ID})
	if msg != nil {
		err = r.repo.UpdateInMessage(message)
		if err != nil {
			logger.Errorf("[Handle InMsg] Failed to update msg %s, %s", message.ID, err)
			return err
		}

		logger.Infof("[Handle InMsg] Updated msg %s", message.ID)
		return nil
	}

	err = r.repo.AddInMessage(message)
	if err != nil {
		logger.Errorf("[Handle InMsg] Failed to insert msg %s, %s", message.ID, err)
		return err
	}
	logger.Infof("[Handle InMsg] Inserted msg %s", message.ID)
	return nil
}

func (r *inHandler) callAPI(message *InMessage) (*http.Response, error) {
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

func (r *inHandler) getPreviousMessage(message InMessage, routingKey string) (
	*InMessage, error) {

	query := bson.M{
		"origin_model":     message.OriginModel,
		"origin_code":      message.OriginCode,
		"routing_key.name": routingKey,
	}
	prevMsg, err := r.repo.GetSingleInMessage(query)

	if err != nil {
		return nil, err
	}

	return prevMsg, nil
}
