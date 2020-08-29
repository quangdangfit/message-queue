package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/quangdangfit/gosdk/utils/logger"

	"message-queue/app/schema"
	"message-queue/app/services"
	"message-queue/pkg/app"
)

type Routing struct {
	service services.RoutingService
}

func NewRouting(service services.RoutingService) *Routing {
	return &Routing{service: service}
}

// Create Routing Key godoc
// @Tags Routing Keys
// @Summary create routing key
// @Description api create routing key
// @Accept  json
// @Produce json
// @Param Body body schema.RoutingCreateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Header 200 {string} Token "qwerty"
// @Router /api/v1/queue/routing_keys [post]
func (r *Routing) Create(c *gin.Context) {
	var bodyParam schema.RoutingCreateParam
	if err := c.Bind(&bodyParam); err != nil {
		logger.Error("Failed to bind body: ", err)
		app.ResError(c, err, 400)
		return
	}

	validate := validator.New()
	if err := validate.Struct(bodyParam); err != nil {
		logger.Error("Body is invalid: ", err)
		app.ResError(c, err, 400)
		return
	}

	message, err := r.service.Create(c, &bodyParam)
	if err != nil {
		logger.Error("Failed to publish message: ", err)
		app.ResError(c, err, 400)
		return
	}

	app.ResSuccess(c, message)
}
