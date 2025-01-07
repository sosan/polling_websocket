package middlewares

import (
	"net/http"
	"polling_websocket/pkg/domain/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateOnGetWorkflow() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: validations
		ctx.Next()
	}
}

func ValidateUserAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func ValidateGetGoogleSheet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actionID := ctx.Param("idaction")
		userID := ctx.Param("iduser")

		// TODO: better validation for actionID and userID
		if strings.TrimSpace(actionID) == "" {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		if strings.TrimSpace(userID) == "" {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
