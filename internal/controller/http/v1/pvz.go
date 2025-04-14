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
		routes.GET("/", middleware.Role(entity.Moderator, entity.Employee), r.getFilterList)

		routes.POST("/", middleware.Role(entity.Moderator), r.create)
		routes.POST("/:id/close_last_reception", middleware.Role(entity.Employee), r.closeLastReception)
		routes.POST("/:id/delete_last_product", middleware.Role(entity.Employee), r.deleteLastProduct)
	}
}

func (r *pvzRoutes) getFilterList(ctx *gin.Context) {
	params := generated.GetPvzParams{}
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
		return
	}

	if params.Page == nil || *params.Page < 1 {
		defaultPage := 1
		params.Page = &defaultPage
	}

	if params.Limit == nil || *params.Limit < 0 {
		defaultLimit := 10
		params.Limit = &defaultLimit
	}

	if params.StartDate != nil {
		if params.StartDate.IsZero() {
			params.StartDate = nil
		}
	}

	if params.EndDate != nil {
		if params.EndDate.IsZero() {
			params.EndDate = nil
		}
	}

	r.l.Info("params value", params)

	pvzList, err := r.pvzUC.ListByReceptionDate(
		ctx.Request.Context(),
		params.StartDate,
		params.EndDate,
		*params.Page,
		*params.Limit,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
		return
	}

	type receptionResponse struct {
		Reception entity.Reception `json:"reception"`
		Products  []entity.Product `json:"products"`
	}

	type pvzResponse struct {
		Pvz        entity.Pvz          `json:"pvz"`
		Receptions []receptionResponse `json:"receptions"`
	}

	response := make([]pvzResponse, 0)

	for _, pvz := range pvzList {
		receptions, err := r.receptionUC.GetByPvzId(ctx.Request.Context(), pvz.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
			return
		}

		receptionsResp := make([]receptionResponse, 0)

		for _, rec := range receptions {
			products, err := r.productUC.GetByReceptionId(ctx.Request.Context(), rec.ID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
				return
			}

			receptionsResp = append(receptionsResp, receptionResponse{
				Reception: rec,
				Products:  products,
			})
		}

		response = append(response, pvzResponse{
			Pvz:        pvz,
			Receptions: receptionsResp,
		})

	}

	ctx.JSON(http.StatusOK, response)
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
