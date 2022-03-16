package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sendya/pkg/ginx/resp"

	"example/api/types"
	"example/internal/service"
)

type UserController struct {
	userService *service.UserService
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

	user, err := c.userService.Login(data.Username, data.Password)
	if err != nil {
		resp.BadRequest(ctx, err)
		return
	}

	// mark security data
	user.Salt = ""
	user.Password = ""

	resp.JSON(ctx, user)
}
