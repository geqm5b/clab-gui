package handlers

import (
	"clab-gui/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LabActionRequest struct {
    Name string `json:"name"`
}


func GetLabsHandler(c *gin.Context) {
	// path de los laboratorios
	const labsPath = "./clab-labs"
	labs, err := service.GetLabs(labsPath)
	if err != nil {
		// Devolvemos un error 500 al navegador
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el directorio de labs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"labs": labs})
}

func DeployLabHandler(c *gin.Context) {
	var request LabActionRequest
	const labsPath = "./clab-labs"
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	if err := service.DeployLab(request.Name, labsPath); err != nil {
        // Si falla, avisamos al cliente con un 500 y el mensaje real
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo el deploy: " + err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "message": "Desplegando " + request.Name,
        "status": "ok",
    })
}
