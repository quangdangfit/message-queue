package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/quangdangfit/gosdk/utils/logger"

	"gomq/app/queue"
	"gomq/app/services"
	"gomq/utils"
)

const (
	ResendOutMessageLimit = 100
	RetryInMessageLimit   = 100
)

type Cron struct {
	inService services.InMessageService
	pub       queue.Publisher
}

func NewCron(service services.InMessageService, publisher queue.Publisher) *Cron {
	return &Cron{
		inService: service,
		pub:       publisher,
	}
}

func (cron *Cron) Resend(c echo.Context) (err error) {
	logger.Info("Start cronjob resend wait messages")

	go cron.pub.CronResend(ResendOutMessageLimit)
	return c.JSON(http.StatusOK, utils.MsgResponse(utils.StatusOK, nil))
}

//func (s *Cron) Retry(c *echo.Context) thttp.Response {
//	logger.Info("Start cronjob resend wait messages")
//
//	go s.inService.CronRetry(RetryInMessageLimit)
//	return thttp.Response{}
//}
//
//func (s *Cron) RetryPrevious(c *echo.Context) thttp.Response {
//	logger.Info("Start cronjob resend wait previous messages")
//
//	go s.inService.CronRetryPrevious(RetryInMessageLimit)
//	return thttp.Response{}
//}
