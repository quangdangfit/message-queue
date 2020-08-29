package impl

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/logger"

	"message-queue/app/models"
	"message-queue/app/queue"
	"message-queue/app/repositories"
	"message-queue/app/schema"
	"message-queue/app/services"
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

func (o *outService) CronResend(limit int) error {
	query := schema.OutMsgQueryParam{
		Status: models.OutMessageStatusWait,
	}
	messages, _ := o.repo.List(&query, limit)
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
