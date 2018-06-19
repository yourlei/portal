package app

import (
	"net/http"
	"encoding/json"

	"portal/util"
	"portal/model"
	"portal/service"
	"github.com/gin-gonic/gin"
)

// Create app
func CreateApp(c *gin.Context) {
	type App struct {
		Name string `json:"name,omitempty"`
	}
	var jsonBody App
	err := c.BindJSON(&jsonBody)
	if err != nil {
		util.RespondBadRequest(c)
		return
	}
	code, msg := service.CreateApp(jsonBody.Name)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"error": gin.H{
			"msg": msg,
		},
	})
}
// Get List
func GetAppList(c *gin.Context) {
	var (
		queryJson *model.AppQueryBody
		code int = 0
	)
	// json string 转为 struct
	if err := json.Unmarshal([]byte(c.Query("query")), &queryJson); err != nil {
		util.RespondBadRequest(c)
		return
	}
	// check time range
	where := queryJson.Where
	if util.CompareDate(where.CreatedAt.Gt, where.CreatedAt.Lt) ||
	   util.CompareDate(where.UpdatedAt.Gt, where.UpdatedAt.Lt) {
		util.RespondBadRequest(c)
		return	
	}
	res, msg := service.GetAppList(queryJson)
	if msg != nil {
		code = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"error": gin.H{
			"msg": msg,
		},
		"data": res,
		"total": len(res),
	})
}