package controller

import (
	"example/api/types"

	"github.com/gin-gonic/gin"
	"github.com/sendya/pkg/ginx/resp"
)

type UserController struct {
}

func NewUserController(v1 *gin.RouterGroup) {
	c := &UserController{}

	user := v1.Group("/user")
	{
		user.POST("/login", c.login)
	}
}

func (c *UserController) login(ctx *gin.Context) {
	var data types.UserRequest
	if err := ctx.ShouldBindJSON(&data); resp.BadRequest(ctx, err) {
		return
	}

	resp.OK(ctx)
}
