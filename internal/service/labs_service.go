package service

import (
	"log"
	"os"
	"strings"
    "os/exec"
    "path/filepath"
	"fmt"
)

// ---La Estructura de Datos ---
type Lab struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// --- Funcion buscar Labs ---
//
// Recibe el path de los labs como argumento y retorna una lista con los labs,
// encontrados.
func GetLabs(basePath string) ([]Lab, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		log.Printf("Error al leer el directorio %s: %v", basePath, err)
		return nil, err
	}
	var labs []Lab
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".clab.yml") {
			labs = append(labs, Lab{Name: file.Name()})
		}
	}
	return labs, nil
}
// --- Funcion desplegar Labs ---
//
// Recibe el path de los labs como argumento y el nombre del lab.
func DeployLab(labName string, basePath string) (error) {
	fullPath := filepath.Join(basePath, labName)
	// revisar existencia del lab
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return fmt.Errorf("el archivo %s no existe", labName)
    }
	// Preparar el comando
	cmd := exec.Command("containerlab", "deploy", "-t", fullPath)
	// Configurar salida para ver logs en consola
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run() 
}


// DEPLOY Y DESTROY SON FUNCIONES MUY SIMILARES, VER SI SE PUEDE HACER UNA SOLA FUNC,
//RECIBIENDO LA ORDEN POR PARAMETRO Y USANDO UN CASE/SWITCH PARA ELEGIR LA ACCION.
func DestroyLab(labName string, basePath string) (error) {
	fullPath := filepath.Join(basePath, labName)
	// revisar existencia del lab
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return fmt.Errorf("el archivo %s no existe", labName)
    }
	// Preparar el comando
	cmd := exec.Command("containerlab", "destroy", "-t", fullPath)
	// Configurar salida para ver logs en consola
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run() 
}


//AGREGAR UNA FUNCION QUE REVISE EL ESTADO DEL LAB 
//