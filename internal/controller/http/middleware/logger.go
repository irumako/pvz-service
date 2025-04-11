package middleware

import (
	"github.com/gin-gonic/gin"
	"pvz-service/pkg/logger"
	"strconv"
	"strings"
	"time"
)

func buildRequestMessage(ctx *gin.Context, duration time.Duration) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.URL.Path)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(ctx.Writer.Status()))
	result.WriteString(" - ")
	result.WriteString(duration.String())

	return result.String()
}

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		l.Info(buildRequestMessage(ctx, duration))
	}
}
