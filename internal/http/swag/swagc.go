//go:build doc
// +build doc

package swag

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example/third_party/swagger"
)

func InitSwagger(r *gin.Engine) {
	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("doc.json"),
			// ginSwagger.DefaultModelsExpandDepth(-1),
			ginSwagger.DocExpansion("none"),
		),
	)
}
