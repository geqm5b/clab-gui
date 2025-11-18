package handlers

import (
	"clab-gui/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLabsHandler(c *gin.Context) {
	const labsPath = "/opt/containerlab/labs"
	labs, err := service.GetLabs(labsPath)
	if err != nil {
		// Devolvemos un error 500 al navegador
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el directorio de labs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"labs": labs})
}
