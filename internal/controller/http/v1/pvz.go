package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	generated "pvz-service/internal/controller/http/dto"
	"pvz-service/internal/controller/http/middleware"
	"pvz-service/internal/entity"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/logger"
)

type pvzRoutes struct {
	pvzUC       usecase.Pvz
	receptionUC usecase.Reception
	productUC   usecase.Product
	l           logger.Interface
}

func NewPvzRoutes(
	apiV1 *gin.RouterGroup,
	pvzUC usecase.Pvz,
	receptionUC usecase.Reception,
	productUC usecase.Product,
	l logger.Interface,
) {
	r := &pvzRoutes{pvzUC, receptionUC, productUC, l}

	routes := apiV1.Group("/pvz")
	{
		routes.GET("/", r.getFilterList)

		routes.POST("/", middleware.Role(entity.Moderator), r.create)
		routes.POST("/:id/close_last_reception", r.closeLastReception)
		routes.POST("/:id/delete_last_product", middleware.Role(entity.Employee), r.deleteLastProduct)
	}
}

func (r *pvzRoutes) getFilterList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *pvzRoutes) create(ctx *gin.Context) {
	body := generated.PostPvzJSONRequestBody{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
		return
	}

	switch body.City {
	case generated.Moscow:
	case generated.Kazan:
	case generated.SaintPetersburg:
	default:
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: "В этом городе нельзя открыть ПВЗ"})
		return
	}

	bodyPvz := entity.Pvz{ID: body.Id.String(), RegistrationDate: *body.RegistrationDate, City: string(body.City)}
	pvz, err := r.pvzUC.Create(ctx.Request.Context(), bodyPvz)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный запрос"})
		return
	}

	ctx.JSON(http.StatusCreated, pvz)
}

func (r *pvzRoutes) closeLastReception(ctx *gin.Context) {
	idParam := ctx.Param("id")
	pvzId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный формат UUID"})
		return
	}

	reception, err := r.receptionUC.CloseLastOpenOnPvz(ctx.Request.Context(), pvzId.String())
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoOpenReception):
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Приемка уже закрыта"})
		default:
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный запрос"})
		}
		return
	}

	ctx.JSON(http.StatusOK, reception)
}

func (r *pvzRoutes) deleteLastProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	pvzId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный формат UUID"})
		return
	}

	err = r.productUC.Delete(ctx.Request.Context(), pvzId.String())
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoOpenReception):
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Приемка уже закрыта"})
		case errors.Is(err, entity.ErrNoProductsToDelete):
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Нет товаров для удаления"})
		default:
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный запрос"})
		}
		return
	}

	ctx.Status(http.StatusOK)
}
