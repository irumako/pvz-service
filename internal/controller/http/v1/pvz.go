package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type pvzRoutes struct {
	u usecase.Pvz
	l logger.Interface
}

func NewPvzRoutes(apiV1 *gin.RouterGroup, u usecase.Pvz, l logger.Interface) {
	r := &pvzRoutes{u, l}

	routes := apiV1.Group("/pvz")
	{
		routes.GET("/", r.getFilterList)

		routes.POST("/", r.create)
		routes.POST("/:id/close_last_reception", r.closeLastReception)
		routes.POST("/:id/delete_last_product", r.deleteLastProduct)
	}
}

func (r *pvzRoutes) getFilterList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *pvzRoutes) create(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *pvzRoutes) closeLastReception(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *pvzRoutes) deleteLastProduct(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
