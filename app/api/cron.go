package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gosdk/utils/logger"

	"message-queue/app/services"
	"message-queue/pkg/app"
)

const (
	ResendOutMessageLimit = 100
	RetryInMessageLimit   = 100
)

type Cron struct {
	inService  services.InService
	outService services.OutService
}

func NewCron(inService services.InService, outService services.OutService) *Cron {
	return &Cron{
		inService:  inService,
		outService: outService,
	}
}

// Resend godoc
// @Tags Retry
// @Summary api resend failed out messages
// @Description api resend `failed` out messages
// @Success 200 {object} app.Response
// @Router /api/v1/cron/resend [post]
func (cron *Cron) Resend(c *gin.Context) {
	logger.Info("Start cronjob resend wait messages")
	go cron.outService.CronResend(ResendOutMessageLimit)
	app.ResOK(c)
}

// Retry godoc
// @Tags Retry
// @Summary api retry `wait retry` in messages
// @Description api resend `wait retry` in messages, message will change status to
// failed when retry more than 3 times
// @Success 200 {object} app.Response
// @Router /api/v1/cron/retry [post]
func (s *Cron) Retry(c *gin.Context) {
	logger.Info("Start cronjob resend wait messages")
	go s.inService.CronRetry(RetryInMessageLimit)
	app.ResOK(c)
}

// Retry Previous godoc
// @Tags Retry
// @Summary api retry `wait retry previous` in messages
// @Description api resend `wait retry previous` in messages, just retry in messages
// have previous message in status (cancel, success)
// @Success 200 {object} app.Response
// @Router /api/v1/cron/retry_previous [post]
func (s *Cron) RetryPrevious(c *gin.Context) {
	logger.Info("Start cronjob resend wait previous messages")
	go s.inService.CronRetryPrevious(RetryInMessageLimit)
	app.ResOK(c)
}
