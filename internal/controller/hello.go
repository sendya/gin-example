package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloController struct {
}

func NewHelloController(r *gin.Engine) {
	c := &HelloController{}

	r.GET("/*any", c.hello)
}

func (c *HelloController) hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, fmt.Sprintf("hello %s", ctx.Request.URL.Path))
}
