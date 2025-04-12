package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	generated "pvz-service/internal/controller/http/dto"
	"pvz-service/internal/usecase"
	"pvz-service/pkg/jwt"
	"pvz-service/pkg/logger"
)

type userRoutes struct {
	u   usecase.User
	l   logger.Interface
	jwt *jwt.Jwt
}

func NewUserRoutes(apiV1 *gin.RouterGroup, u usecase.User, l logger.Interface, jwt *jwt.Jwt) {
	r := &userRoutes{u, l, jwt}

	routes := apiV1.Group("/")
	{
		routes.POST("/dummyLogin", r.getDummyLogin)
		routes.POST("/register", r.register)
		routes.POST("/login", r.login)
	}
}

func (r *userRoutes) getDummyLogin(ctx *gin.Context) {
	var req generated.PostDummyLoginJSONRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: err.Error()})
		return
	}

	if req.Role != generated.PostDummyLoginJSONBodyRoleEmployee && req.Role != generated.PostDummyLoginJSONBodyRoleModerator {
		ctx.JSON(http.StatusBadRequest, generated.Error{Message: "Неверный запрос"})
		return
	}

	var token generated.Token

	token, err := r.jwt.GenerateToken("dummy-user-id", "dummy@example.com", string(req.Role))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, generated.Error{Message: "Ошибка сервера"})
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (r *userRoutes) register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *userRoutes) login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
