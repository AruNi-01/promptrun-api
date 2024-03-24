package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/service"
	"strconv"
)

func ModelList(c *gin.Context) {
	modelList, e := service.ModelList(c)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(modelList))
}

func FindModelById(c *gin.Context) {
	modelId, _ := strconv.Atoi(c.Param("id"))
	promptModel, e := service.FindModelById(modelId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(promptModel))
}
