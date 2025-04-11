package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type productRoutes struct {
	u usecase.Product
	l logger.Interface
}

func NewProductRoutes(apiV1 *gin.RouterGroup, u usecase.Product, l logger.Interface) {
	r := &productRoutes{u, l}

	routes := apiV1.Group("/products")
	{
		routes.POST("/", r.add)
	}
}

func (r *productRoutes) add(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
