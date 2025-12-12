package service

import (
	"clab-gui/internal/converter" // importar paquete converter
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/goccy/go-yaml"
)

type Lab struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// yml parsing strucs
type ymlNode struct {
	Kind string `yaml:"kind"`
}
type ymlTopology struct {
	Nodes map[string]ymlNode `yaml:"nodes"`
}
type clabConfig struct {
	Topology ymlTopology `yaml:"topology"`
}


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

// funcion para desplegar un Lab
func DeployLab(labName string, basePath string) error {
	fullPath := filepath.Join(basePath, labName)

	// Verificar que el archivo exista
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo %s no existe", fullPath)
	}

	// Leer el YAML
	clabLab, err := readYml(fullPath)
	if err != nil {
		return err
	}

	// Obtener los bridges
	bridgesList, err := getBridgesNames(clabLab)
	if err != nil {
		return err
	}

	// Crear bridges si no existen
	for _, bridge := range bridgesList {
		if err := createBridge(bridge); err != nil {
			return fmt.Errorf("error creando bridge %s: %w", bridge, err)
		}
	}

	// Deploy containerlab
	cmd := exec.Command("containerlab", "deploy", "-t", fullPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// funcion para destruir un Lab
func DestroyLab(labName string, basePath string) error {
	fullPath := filepath.Join(basePath, labName)

	// Verificar que el archivo exista
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo %s no existe", fullPath)
	}

	// Leer el YAML antes de destruir el lab
	clabLab, err := readYml(fullPath)
	if err != nil {
		return fmt.Errorf("error leyendo %s: %w", fullPath, err)
	}

	// Obtener bridges definidos en el YAML
	bridgesList, err := getBridgesNames(clabLab)
	if err != nil {
		return fmt.Errorf("error obteniendo bridges: %w", err)
	}

	// Ejecutar containerlab destroy
	cmd := exec.Command("containerlab", "destroy", "-t", fullPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error destruyendo lab con containerlab: %w", err)
	}

	// Borrar bridges después del destroy
	for _, bridge := range bridgesList {
		if err := deleteBridge(bridge); err != nil {
			return fmt.Errorf("error eliminando bridge %s: %w", bridge, err)
		}
	}

	return nil
}

func CreateLabFromStream(fileStream io.Reader, originalFilename string, destDir string) (string, error) {
	// Determinar nombre del Lab basado en el archivo
	baseName := strings.TrimSuffix(originalFilename, filepath.Ext(originalFilename))
	// Limpieza del onmbre
	if idx := strings.LastIndex(baseName, "."); idx != -1 {
		baseName = baseName[:idx] // Quedarse solo con el nombre base
	}
	
	// Nombre final del archivo YAML
	finalName := baseName + ".clab.yml"
	finalPath := filepath.Join(destDir, finalName)

	// Convertir XML a YAML en memoria 
	yamlBytes, err := converter.ConvertDrawioToYaml(fileStream, baseName)
	if err != nil {
		return "", fmt.Errorf("falló la conversión del diseño: %v", err)
	}


	// Asegurar que existe la carpeta labs
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creando directorio de labs: %v", err)
	}

	// Escribir el resultado (*.clab.yml) en el disco
	if err := os.WriteFile(finalPath, yamlBytes, 0644); err != nil {
		return "", fmt.Errorf("error escribiendo el archivo YAML: %v", err)
	}

	return finalName, nil
}

//lee un arvhico clab.yml del sistema de archivos del os y lo devulve
func readYml(clabLabPath string) ([]byte, error) {
	// os.ReadFile recibe la ruta del archivo y devuelve el contenido ([]byte) o un error.
	contenido, err := os.ReadFile(clabLabPath)

	if err != nil {

		return nil, fmt.Errorf("error al leer el archivo %s: %w", clabLabPath, err)
	}

	return contenido, nil
}

// parsea el yml y devuelve una lista con los nombres de los bridges
func getBridgesNames(clabLab []byte) ([]string, error) {
	var config clabConfig

	//  mapeo a la estructura 'config'
	if err := yaml.Unmarshal(clabLab, &config); err != nil {
		return nil, fmt.Errorf("error al deserializar YAML: %w", err)
	}

	var bridges []string

	// loop  sobre el mapa de nodos
	for bridgeName, nodeData := range config.Topology.Nodes {
		if nodeData.Kind == "bridge" {
			bridges = append(bridges, bridgeName)
		}
	}
	return bridges, nil
}

//recibe un nombre de bridge y lo crea en el OS. ej. ip link add tesoyunbride type bridge
func createBridge(bridge string) error {
    exists, err := bridgeExists(bridge)
    if err != nil {
        return fmt.Errorf("error verificando existencia del bridge %s: %w", bridge, err)
    }
    if exists {
        return fmt.Errorf("bridge '%s' ya existe", bridge)
    }

    // Crear bridge
    cmd := exec.Command("ip", "link", "add", bridge, "type", "bridge")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("error creando bridge %s: %w", bridge, err)
    }

    // Levantar bridge
    cmd = exec.Command("ip", "link", "set", bridge, "up")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("error levantando bridge %s: %w", bridge, err)
    }

    return nil
}


func deleteBridge(bridge string) error {
    exists, err := bridgeExists(bridge)
    if err != nil {
        return err
    }

    if !exists {
        return nil
    }

    // Bajar bridge antes de eliminarlo 
    cmdDown := exec.Command("ip", "link", "set", bridge, "down")
    cmdDown.Stdout = os.Stdout
    cmdDown.Stderr = os.Stderr
    if err := cmdDown.Run(); err != nil {
        return fmt.Errorf("error bajando bridge %s: %w", bridge, err)
    }

    // Borrar el bridge
    cmd := exec.Command("ip", "link", "delete", bridge)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("error eliminando bridge %s: %w", bridge, err)
    }

    return nil
}


func bridgeExists(name string) (bool, error) {
	cmd := exec.Command("sh", "-c", "ip -o link show type bridge | awk -F': ' '{print $2}' | cut -d'@' -f1")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	bridges := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, b := range bridges {
		if b == name {
			return true, nil
		}
	}
	return false, nil
}
