package services

type InService interface {
	Consume()
	CronRetry(limit int) error
	CronRetryPrevious(limit int) error
}
