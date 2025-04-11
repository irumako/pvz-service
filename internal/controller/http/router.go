package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/internal/controller/http/middleware"
	v1 "pvz-service/internal/controller/http/v1"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type Usecases struct {
	PvzUC       usecase.Pvz
	ReceptionUC usecase.Reception
	UserUC      usecase.User
	ProductUC   usecase.Product
}

func SetRouters(
	r *gin.Engine,
	l logger.Interface,
	uc Usecases,
) {
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(l))

	r.GET("/healthcheck", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	routes := r.Group("/api/v1")
	//v1.NewUserRoutes(routes, uc.GetUserNameList)
	//routes.Use(middleware.BasicAuthMiddleware(uc.GetUser))
	{
		v1.NewPvzRoutes(routes, uc.PvzUC, l)
		v1.NewUserRoutes(routes, uc.UserUC, l)
		v1.NewReceptionRoutes(routes, uc.ReceptionUC, l)
		v1.NewProductRoutes(routes, uc.ProductUC, l)
	}
}
