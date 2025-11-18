package main

import (
	"log"

	"clab-gui/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// --- Servir el Frontend (el "sitio" en web/) ---
	// 1. Sirve los archivos estáticos (app.js) desde la ruta /static
	router.Static("/static", "./web")

	// 2. Sirve el index.html en la raíz "/"
	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// --- Servir el Backend (API) ---
	api := router.Group("/api")
	{
		api.GET("/labs", handlers.GetLabsHandler)

	}

	log.Println("Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}
