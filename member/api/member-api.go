package api

import (
	"github.com/gin-gonic/gin"
	"go-sample/common"
	"go-sample/member/app"
	"go-sample/member/query"
	"net/http"
	"strconv"
)

func CreateMemberApi(
	queryService query.MemberQueryService,
	registerService app.RegisterService,
	updateService app.UpdateService,
) *MemberApi {
	return &MemberApi{
		queryService, registerService, updateService,
	}
}

type MemberApi struct {
	query.MemberQueryService
	app.RegisterService
	app.UpdateService
}

// HandlerRegisterer
func (api *MemberApi) RegisterHandler(router *gin.Engine) {
	router.GET("/members", api.MemberList)
	router.GET("/members/:id", api.Member)

	router.POST("/members", api.RegisterMember)
	router.PUT("/members", api.UpdateMember)
}

func (api *MemberApi) MemberList(c *gin.Context) {
	members, err := api.MemberQueryService.GetMembers()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, members)
}

func (api *MemberApi) Member(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	member, err := api.MemberQueryService.GetMember(id)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, member)
}

func (api *MemberApi) RegisterMember(c *gin.Context) {
	var req app.MemberRegistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	id, err := api.Register(req)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": id})
}

func (api *MemberApi) UpdateMember(c *gin.Context) {
	var req app.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	err := api.UpdateService.UpdateMember(req)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	c.String(http.StatusOK, "")
}
