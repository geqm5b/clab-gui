package service

import (
	"log"
	"os"
	"strings"
)

// ---La Estructura de Datos ---
//
// Devolvemos una lista de 'Lab', no solo una lista de 'string'.
type Lab struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// --- Funcion buscar Labs ---
//
// Recibe el path de los labs como argumento y retorna una lista con los labs,
// encontrados.
func GetLabs(path string) ([]Lab, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		log.Printf("Error al leer el directorio %s: %v", path, err)
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

//func DeployLab(path string) ([]Lab, error) {
//
//}
