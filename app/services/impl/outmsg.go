package impl

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/queue"
	"message-queue/app/repositories"
	"message-queue/app/schema"
	"message-queue/app/services"
)

const (
	ResendOutMessageLimit = 100
)

type outService struct {
	pub  queue.Publisher
	repo repositories.OutRepository
}

func NewOutService(pub queue.Publisher, repo repositories.OutRepository) services.OutService {
	return &outService{
		pub:  pub,
		repo: repo,
	}
}

func (o *outService) List(ctx context.Context, query *schema.OutMsgQueryParam) (*[]models.OutMessage, *paging.Paging, error) {
	rs, pageInfo, err := o.repo.List(query)
	if err != nil {
		logger.Errorf("Failed to get list out messages")
		return nil, nil, err
	}
	return rs, pageInfo, nil
}

func (o *outService) Update(ctx context.Context, id string, body *schema.OutMsgUpdateParam) (*models.OutMessage, error) {
	rs, err := o.repo.Update(id, body)
	if err != nil {
		logger.Errorf("Failed to update out message %s, error: ", id)
		return nil, err
	}
	return rs, nil
}

func (o *outService) Publish(ctx context.Context, message *models.OutMessage) error {
	err := o.pub.Publish(message, true)
	if err != nil {
		logger.Errorf("Failed to publish msg %s, %s", message.ID, err)
	}

	err = o.repo.Create(message)
	if err != nil {
		logger.Errorf("Failed to create out msg %s", message.ID)
		return err
	}
	return nil
}

func (o *outService) CronResend() error {
	query := schema.OutMsgQueryParam{
		Status: models.OutMessageStatusWait,
		Limit:  ResendOutMessageLimit,
	}
	messages, _, _ := o.repo.List(&query)
	if messages == nil {
		logger.Info("[Resend Message] Not found any wait message!")
		return nil
	}

	logger.Infof("[Resend Message] Found %d wait messages!", len(*messages))
	for _, msg := range *messages {
		o.pub.Publish(&msg, true)
	}
	logger.Info("[Resend Message] Finish!")

	return nil
}
