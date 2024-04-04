package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
	"strconv"
)

func PromptList(c *gin.Context) {
	var promptListReq service.PromptListReq
	if err := c.ShouldBindJSON(&promptListReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		prompts, e := promptListReq.PromptList(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(
			&service.PromptListResp{
				Prompts: prompts,
				Rows:    promptListReq.Paginate.Rows,
			},
		))
	}
}

func FindById(c *gin.Context) {
	promptId, _ := strconv.Atoi(c.Param("id"))
	prompt, e := service.FindPromptById(c, promptId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(prompt))
}

func FindFullInfoById(c *gin.Context) {
	promptId, _ := strconv.Atoi(c.Param("id"))
	promptDetail, e := service.FindPromptFullInfoById(c, promptId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(promptDetail))
}

func FindImgListByPromptId(c *gin.Context) {
	promptId, _ := strconv.Atoi(c.Param("promptId"))
	promptImgList, e := service.FindPromptImgListByPromptId(c, promptId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(promptImgList))
}

func FindMasterImgByPromptId(c *gin.Context) {
	promptId, _ := strconv.Atoi(c.Param("promptId"))
	promptImg, e := service.FindPromptMasterImgByPromptId(promptId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(promptImg))
}

func FindMasterImgListByPromptIds(c *gin.Context) {
	var promptMasterImgService service.PromptMasterImgListReq
	if err := c.ShouldBindJSON(&promptMasterImgService); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		promptImgList, e := promptMasterImgService.FindMasterImgListByPromptIds(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(promptImgList))
	}
}

func FindListByBuyerId(c *gin.Context) {
	var promptListByBuyerIdReq service.PromptListByBuyerIdReq
	if err := c.ShouldBindJSON(&promptListByBuyerIdReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		prompts, e := promptListByBuyerIdReq.FindListByBuyerId(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(
			&service.PromptListResp{
				Prompts: prompts,
				Rows:    promptListByBuyerIdReq.Paginate.Rows,
			},
		))
	}
}
