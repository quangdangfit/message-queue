package grpc

import (
	"context"
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/validator"

	"message-queue/app/models"
	"message-queue/app/schema"
	"message-queue/app/services"
)

type OutRPC struct {
	service services.OutService
}

func NewOutRPC(service services.OutService) *OutRPC {
	return &OutRPC{service: service}
}

func (o *OutRPC) Publish(body map[string]interface{}, reply *string) error {
	data, err := json.Marshal(body)
	if err != nil {
		*reply = err.Error()
		return err
	}

	var bodyParam schema.OutMsgCreateParam
	err = json.Unmarshal(data, &bodyParam)
	if err != nil {
		*reply = err.Error()
		return err
	}

	validate := validator.New()
	if err := validate.Validate(bodyParam); err != nil {
		err = errors.Wrap(err, "body is invalid")
		*reply = err.Error()
		logger.Error(err)
		return err
	}

	message, err := o.prepareMessage(bodyParam)
	if err != nil {
		err = errors.Wrap(err, "failed to parse out message")
		*reply = err.Error()
		logger.Error(err)
		return err
	}

	ctx := context.Background()
	err = o.service.Publish(ctx, message)
	if err != nil {
		err = errors.Wrap(err, "failed to publish message")
		*reply = err.Error()
		logger.Error(err)
		return err
	}

	return nil
}

func (o *OutRPC) prepareMessage(bodyParam schema.OutMsgCreateParam) (*models.OutMessage, error) {
	message := models.OutMessage{}
	err := copier.Copy(&message, &bodyParam)

	if err != nil {
		return &message, err
	}
	message.Status = models.OutMessageStatusWait

	return &message, nil
}
