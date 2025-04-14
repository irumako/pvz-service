package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	generated "pvz-service/internal/controller/http/dto"
	"pvz-service/internal/controller/http/middleware"
	"pvz-service/internal/entity"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type receptionRoutes struct {
	receptionUC usecase.Reception
	l           logger.Interface
}

func NewReceptionRoutes(apiV1 *gin.RouterGroup, receptionUC usecase.Reception, l logger.Interface) {
	r := &receptionRoutes{receptionUC, l}

	routes := apiV1.Group("/receptions")
	{
		routes.POST("/", middleware.Role(entity.Employee), r.create)
	}
}

func (r *receptionRoutes) create(ctx *gin.Context) {
	body := generated.PostReceptionsJSONRequestBody{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
		return
	}

	reception, err := r.receptionUC.Create(ctx.Request.Context(), body.PvzId.String())
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUnclosedReception):
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Есть незакрытая приемка"})
		default:
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный запрос"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, reception)
}
