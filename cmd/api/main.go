package main

import (
	"log"

	"clab-gui/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Aumentar límite de memoria para subidas 
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Servir el Frontend
	router.Static("/static", "./web")
	router.Static("/drawio", "./drawio_static")

	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// Servir el Backend (API)
	api := router.Group("/api")
	{
		api.GET("/getLabs", handlers.GetLabsHandler)
		api.POST("/deployLab", handlers.DeployLabHandler)
		api.POST("/destroyLab", handlers.DestroyLabHandler)
		api.POST("/upload", handlers.UploadHandler)
	}

	router.GET("/editor", func(c *gin.Context) {
		c.File("./web/editor.html")
	})

	log.Println("Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}