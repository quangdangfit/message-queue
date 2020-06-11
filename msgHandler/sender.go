package msgHandler

import (
	"gomq/dbs"
	"gomq/models"
	"gomq/msgQueue"
	"time"

	"github.com/jinzhu/copier"
)

type Sender interface {
	ParseMessage(msgRequest models.MessageRequest) (*models.OutMessage, error)
	Send(msgRequest models.MessageRequest) error
}

type sender struct {
	publisher msgQueue.Publisher
}

func NewSender() Sender {
	s := sender{
		publisher: msgQueue.NewPublisher(),
	}
	return &s
}

func (s *sender) ParseMessage(msgRequest models.MessageRequest) (*models.OutMessage, error) {
	message := models.OutMessage{}
	err := copier.Copy(&message, &msgRequest)

	if err != nil {
		return &message, err
	}

	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.Status = dbs.OutMessageStatusWait

	return &message, nil
}

func (s *sender) Send(msgRequest models.MessageRequest) error {
	message, err := s.ParseMessage(msgRequest)

	if err != nil {
		return err
	}

	s.publisher.Publish(message, true)

	dbs.AddOutMessage(message)
	return nil
}
