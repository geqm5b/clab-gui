package main

import (
	"log"

	"clab-gui/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// --- Servir el Frontend---
	//  Sirve los archivos estáticos (app.js) desde la ruta /static
	router.Static("/static", "./web")

	//  Sirve el index.html en la raíz "/"
	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// --- Servir el Backend (API) ---
	api := router.Group("/api")
	{
		api.GET("/getLabs", handlers.GetLabsHandler)
		api.POST("/deployLab", handlers.DeployLabHandler)
		api.POST("/destroyLab", handlers.DestroyLabHandler)
	}

	log.Println("Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}
