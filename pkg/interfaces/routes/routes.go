package routes

import (
	"net/http"
	"polling_websocket/pkg/dimodel"
	"polling_websocket/pkg/interfaces/middlewares"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine, dependencies *dimodel.Dependencies) {
	app.NoRoute(ErrRouter)

	// Routes in groups
	api := app.Group("/api/v1")
	{
		// api.GET("/ping", common.Ping)

		actions := api.Group("/polling")
		{
			actions.GET("/google", dependencies.PollingController.Ping)
			actions.GET("/google/sheets/:iduser/:idaction", middlewares.ValidateGetGoogleSheet(), dependencies.PollingController.GetGoogleSheetByID)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}