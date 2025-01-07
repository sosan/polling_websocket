package controllers

import (
	"net/http"
	"polling_websocket/pkg/domain/models"
	"polling_websocket/pkg/domain/repos"

	"github.com/gin-gonic/gin"
)

type PollingController struct {
	pollingService repos.PollingService
}

func NewPollingController(newPollingService repos.PollingService) *PollingController {
	return &PollingController{pollingService: newPollingService}
}

func (a *PollingController) Ping(ctx *gin.Context) {
	ob := gin.H{
		"test": "test",
	}
	ctx.JSON(200, ob)
}

func (a *PollingController) GetGoogleSheetByID(ctx *gin.Context) {
	actionID := ctx.Param("idaction")
	userID := ctx.Param("iduser")
	data, err := a.pollingService.GetContentGoogleSheetByID(&actionID, &userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseGetGoogleSheetByID{
			Error:  "not generated",
			Status: http.StatusInternalServerError,
		})
		return
	}
	if data == nil { // case empty data
		ctx.JSON(http.StatusOK, models.ResponseGetGoogleSheetByID{
			Error:  "",
			Status: http.StatusOK,
			Data:   "",
		})
		return
	}
	if *data == "" { // case empty data
		ctx.JSON(http.StatusOK, models.ResponseGetGoogleSheetByID{
			Error:  "",
			Status: http.StatusOK,
			Data:   "",
		})
		return
	}
	ctx.JSON(http.StatusCreated, models.ResponseGetGoogleSheetByID{
		Status: http.StatusOK,
		Error:  "",
		Data:   *data,
	})
}
