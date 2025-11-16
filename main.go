package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// Endpoint para la página principal
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Containerlab GUI",
		})
	})

	// Endpoint para la prueba de HTMX
	router.POST("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "<p>¡Respuesta desde el Backend de Go!!</p>")
	})

	router.Run(":8080")
}
