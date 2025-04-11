package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type userRoutes struct {
	u usecase.User
	l logger.Interface
}

func NewUserRoutes(apiV1 *gin.RouterGroup, u usecase.User, l logger.Interface) {
	r := &userRoutes{u, l}

	routes := apiV1.Group("/")
	{
		routes.POST("/dummyLogin", r.getDummyLogin)
		routes.POST("/register", r.register)
		routes.POST("/login", r.login)
	}
}

func (r *userRoutes) getDummyLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *userRoutes) register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *userRoutes) login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
