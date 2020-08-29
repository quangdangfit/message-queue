package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"message-queue/docs"
)

func RegisterDocs(e *gin.Engine) {
	// Swagger info
	docs.SwaggerInfo.Title = "Message Queue"
	docs.SwaggerInfo.Description = "Swagger message queue API document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
