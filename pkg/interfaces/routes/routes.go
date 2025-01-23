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
		actions := api.Group("/polling")
		{
			actions.GET("/ping", dependencies.PollingController.Ping)
			// both uris uses same logic, btw i dont know if notion action will change base logic
			actions.GET("/google/sheets/:iduser/:idaction", middlewares.ValidateRequestAction(), dependencies.PollingController.GetActionByID)
			actions.GET("/notion/:iduser/:idaction", middlewares.ValidateRequestAction(), dependencies.PollingController.GetActionByID)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
