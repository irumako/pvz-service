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

type productRoutes struct {
	productUC usecase.Product
	l         logger.Interface
}

func NewProductRoutes(apiV1 *gin.RouterGroup, productUC usecase.Product, l logger.Interface) {
	r := &productRoutes{productUC, l}

	routes := apiV1.Group("/products")
	{
		routes.POST("/", middleware.Role(entity.Employee), r.add)
	}
}

func (r *productRoutes) add(ctx *gin.Context) {
	body := generated.PostProductsJSONRequestBody{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
		return
	}

	switch body.Type {
	case generated.PostProductsJSONBodyTypeClothes:
	case generated.PostProductsJSONBodyTypeShoes:
	case generated.PostProductsJSONBodyTypeElectronics:
	default:
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный тип товара"})
		return
	}

	product, err := r.productUC.Create(ctx.Request.Context(), body.PvzId.String(), entity.ProductType(body.Type))
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoOpenReception):
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Нет активной приемки"})
		default:
			ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный запрос"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, product)
}
