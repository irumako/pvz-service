package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type receptionRoutes struct {
	u usecase.Reception
	l logger.Interface
}

func NewReceptionRoutes(apiV1 *gin.RouterGroup, u usecase.Reception, l logger.Interface) {
	r := &receptionRoutes{u, l}

	routes := apiV1.Group("/receptions")
	{
		routes.POST("/", r.create)
	}
}

func (r *receptionRoutes) create(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
