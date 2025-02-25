package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harrychenshoplazza/service_tree/internal/services"
)

func SetupRouter(serviceHandler *services.ServiceHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("api/v1/servicetree")
	{
		v1.POST("", serviceHandler.CreateService)
	}
	return r
}
