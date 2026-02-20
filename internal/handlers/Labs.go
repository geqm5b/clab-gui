package handlers

import (
	"clab-gui/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LabActionRequest struct {
	Name string `json:"name"`
}

const labsPath = "./clab-labs" 

// Handler para que devuelve la lista de labs
func GetLabsHandler(c *gin.Context) {
	labs, err := service.GetLabs(labsPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer labs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"labs": labs})
}

// handrler del upload del archivo .drawio 
func UploadHandler(c *gin.Context) {
	// Recibir el stream
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se recibió archivo 'file'"})
		return
	}
	defer file.Close()
	// Llamamos a la funcion 
	createdLabName, err := service.CreateLabFromStream(file, header.Filename, labsPath)
	// Menasaje de error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando el diseño: " + err.Error()})
		return
	}
	// Mensaje de exito
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Laboratorio creado correctamente",
		"lab_name": createdLabName,
	})
}

// Handler para desplegar un lab
// Falta la creacion de bridges!!
func DeployLabHandler(c *gin.Context) {
	var request LabActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}
	if err := service.DeployLab(request.Name, labsPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo deploy: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Desplegando " + request.Name, "status": "ok"})
}

// Handler para destruir un lab
func DestroyLabHandler(c *gin.Context) {
	var request LabActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}
	if err := service.DestroyLab(request.Name, labsPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo destroy: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Eliminando " + request.Name, "status": "ok"})
}