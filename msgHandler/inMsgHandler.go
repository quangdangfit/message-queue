package msgHandler

import (
	"encoding/json"
	"gomq/dbs"
	"gomq/models"
)

type Receiver interface {
	HandleMessage(bytesMsg []byte, routingKey string) (
		message *models.InMessage, err error)

	convertMessage(bytesMsg []byte) (*models.InMessage, error)
	storeMessage(message *models.InMessage, routingKey string) (err error)
	callAPI(message *models.InMessage) (err error)
}

type receiver struct{}

func NewReceiver() Receiver {
	r := receiver{}
	return &r
}

func (r *receiver) HandleMessage(bytesMsg []byte, routingKey string) (
	message *models.InMessage, err error) {

	message, _ = r.convertMessage(bytesMsg)
	err = r.callAPI(message)

	if err != nil {
		message.Status = dbs.InMessageStatusWaitRetry
	}
	r.storeMessage(message, routingKey)

	return message, nil
}

func (r *receiver) convertMessage(bytesMsg []byte) (*models.InMessage, error) {
	var message = models.InMessage{}
	json.Unmarshal(bytesMsg, &message)
	return &message, nil
}

func (r *receiver) storeMessage(message *models.InMessage, routingKey string) (
	err error) {

	message, _ = dbs.AddInMessage(message, routingKey)

	return nil
}

func (r *receiver) callAPI(message *models.InMessage) (err error) {
	return nil
}
