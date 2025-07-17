package routers

import (
	"vodka/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/bioskop", controllers.CreateBioskop)

	router.GET("/bioskop", controllers.GetBioskop)

	router.GET("/bioskop/:id", controllers.GetBioskopById)

	router.PUT("/bioskop/:id", controllers.UpdateBioskop)

	router.DELETE("/bioskop/:id", controllers.DeleteBioskop)
	return router
}
