package msgHandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gomq/dbs"
	"gomq/models"
	"gomq/utils"
	"net/http"
	"time"
	"transport/lib/utils/logger"
)

const RequestTimeout = time.Duration(60 * time.Second)

type Receiver interface {
	HandleMessage(message *models.InMessage, routingKey string) (*models.InMessage, error)
	storeMessage(message *models.InMessage) (err error)
	callAPI(message *models.InMessage) (*http.Response, error)
}

type receiver struct{}

func NewReceiver() Receiver {
	r := receiver{}
	return &r
}

func (r *receiver) HandleMessage(message *models.InMessage, routingKey string) (
	*models.InMessage, error) {

	inRoutingKey, err := dbs.GetRoutingKey(routingKey)
	if err != nil {
		message.Status = dbs.InMessageStatusInvalid
		message.Logs = err.Error()
		r.storeMessage(message)
		return message, err
	}
	message.RoutingKey = *inRoutingKey

	res, err := r.callAPI(message)
	if err != nil {
		message.Status = dbs.InMessageStatusWaitRetry
		message.Logs = err.Error()
	} else if res.StatusCode != http.StatusOK {
		message.Status = dbs.InMessageStatusWaitRetry
		message.Logs = utils.ParseError(*res)
	}

	r.storeMessage(message)

	return message, err
}

func (r *receiver) storeMessage(message *models.InMessage) (err error) {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()

	message, _ = dbs.AddInMessage(message)
	return nil
}

func (r *receiver) callAPI(message *models.InMessage) (*http.Response, error) {
	routingKey := message.RoutingKey

	bytesPayload, _ := json.Marshal(message.Payload)
	req, _ := http.NewRequest(
		routingKey.APIMethod, routingKey.APIUrl, bytes.NewBuffer(bytesPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "ahsfishdi"))

	client := http.Client{
		Timeout: RequestTimeout,
	}
	res, err := client.Do(req)

	if err != nil {
		logger.Errorf("Failed to send request to %s, %s", routingKey.APIUrl, err)
		return res, err
	}

	return res, nil
}
