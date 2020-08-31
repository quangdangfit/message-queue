package services

type InService interface {
	Consume()
	CronRetry() error
	CronRetryPrevious() error
}
