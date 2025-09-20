package web

import (
	"github.com/gin-gonic/gin"
)

func SetupEngine() *gin.Engine {
	engine := gin.Default()
	return engine
}
