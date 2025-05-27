package routes

import (
	"people-enricher/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes routes and returns the router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/people", controllers.AddPerson)
	r.GET("/people", controllers.GetPeople)
	r.PUT("/people/:id", controllers.UpdatePerson)
	r.DELETE("/people/:id", controllers.DeletePerson)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
