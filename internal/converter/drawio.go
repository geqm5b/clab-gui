package converter

import (
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/goccy/go-yaml"
)

// struct datos raw de interfaces
type interfaceRaw struct {
	Name   string
	Ip     string
	Id     string
	Parent string
}

// struct datos raw de los nodos
type nodeRaw struct {
	Name       string
	Kind       string
	Image      string
	Id         string
	Interfaces []interfaceRaw
}

// struct datos raw de los links
type edgeRaw struct {
	Source string
	Target string
}

// struct raiz del yaml
type clab struct {
	Name     string   `yaml:"name"`
	Topology topology `yaml:"topology"`
}

type topology struct {
	Kinds map[string]kindConfig `yaml:"kinds,omitempty"`
	Nodes map[string]node       `yaml:"nodes"`
	Links []edge                `yaml:"links"`
}

// configuracion global de kinds
type kindConfig struct {
	NetworkMode string `yaml:"network-mode,omitempty"`
}

// struct sanitizada del nodo
type node struct {
	Name  string            `yaml:"-"`
	Kind  string            `yaml:"kind"`
	Image string            `yaml:"image,omitempty"`
	Binds []string          `yaml:"binds,omitempty"`
	Env   map[string]string `yaml:"env,omitempty"`
	Exec  []string          `yaml:"exec,omitempty"`
}

// struct sanitizada del link
type edge struct {
	Endpoints []string `yaml:"endpoints,flow"`
}

// estructura auxiliar interna
type ifaceInfo struct {
	NodeName  string
	IfaceName string
	Kind      string
}

// estructura para defaults 
type defaultConfig struct {
	Binds []string
	Env   map[string]string
}

// Helper para navegar al root sin repetir código
func getRootElement(doc *etree.Document) *etree.Element {
	//root := doc.FindElement("/mxfile/diagram/mxGraphModel/root")
	mxfile := doc.SelectElement("mxfile")
	if mxfile == nil {
		return nil
	}
	diagram := mxfile.SelectElement("diagram")
	if diagram == nil {
		return nil
	}
	mxGraphModel := diagram.SelectElement("mxGraphModel")
	if mxGraphModel == nil {
		return nil
	}
	return mxGraphModel.SelectElement("root")
}

// helper para formatear endpoint segun tipo
func formatEndpoint(info ifaceInfo) string {
	if info.Kind == "bridge" {
		return fmt.Sprintf("%s:%s-%s", info.NodeName, info.NodeName, info.IfaceName)
	}
	return fmt.Sprintf("%s:%s", info.NodeName, info.IfaceName)
}

// funcion para obtener los nodos raw
func getNodesRaW(doc *etree.Document) ([]nodeRaw, error) {
	var nodes_raw []nodeRaw
	var interfaces_raw []interfaceRaw

	root := getRootElement(doc)
	if root == nil {
		return nil, fmt.Errorf("xml invalido: no se pudo navegar hasta el elemento root")
	}

	// loop para obtener los nodos e interfaces
	for _, object := range root.SelectElements("object") {
		mxCell := object.FindElement("mxCell")
		if mxCell == nil {
			continue
		}

		parent := mxCell.SelectAttrValue("parent", "")

		// si parent es 1 es un nodo
		if parent == "1" {
			nodes_raw = append(nodes_raw, nodeRaw{
				Name:  object.SelectAttrValue("label", ""),
				Kind:  object.SelectAttrValue("kind", ""),
				Image: object.SelectAttrValue("image", ""),
				Id:    object.SelectAttrValue("id", ""),
			})
		}

		// si parent no es 1 es una interfaz
		if parent != "1" {
			interfaces_raw = append(interfaces_raw, interfaceRaw{
				Name:   object.SelectAttrValue("label", ""),
				Ip:     object.SelectAttrValue("ip", ""),
				Id:     object.SelectAttrValue("id", ""),
				Parent: parent,
			})
		}
	}

	// asignar interfaces a los nodos correspondientes
	for i := range nodes_raw {
		for _, iface := range interfaces_raw {
			if iface.Parent == nodes_raw[i].Id {
				nodes_raw[i].Interfaces = append(nodes_raw[i].Interfaces, iface)
			}
		}
	}

	return nodes_raw, nil
}

// funcion para obtener los edges raw
func getEdgesRaW(doc *etree.Document) ([]edgeRaw, error) {
	var edges_raw []edgeRaw

	root := getRootElement(doc)
	if root == nil {
		return nil, fmt.Errorf("xml invalido: no se pudo navegar hasta el elemento root")
	}

	for _, mxCell := range root.SelectElements("mxCell") {
		if edge := mxCell.SelectAttr("edge"); edge != nil && edge.Value == "1" {
			source := mxCell.SelectAttrValue("source", "")
			target := mxCell.SelectAttrValue("target", "")

			if source != "" && target != "" {
				edges_raw = append(edges_raw, edgeRaw{Source: source, Target: target})
			}
		}
	}
	return edges_raw, nil
}

// Funcion Principal 
func ConvertDrawioToYaml(r io.Reader, labName string) ([]byte, error) {

	// valores por defecto
	defaults := defaultConfig{
		Binds: []string{"/tmp/.X11-unix:/tmp/.X11-unix"},
		Env:   map[string]string{"DISPLAY": ":0"},
	}

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(r); err != nil {
		return nil, err
	}

	// obtener datos raw
	nodesRaw, err := getNodesRaW(doc)
	if err != nil {
		return nil, err
	}
	edgesRaw, err := getEdgesRaW(doc)
	if err != nil {
		return nil, err
	}

	// inicializar topologia
	lab := clab{
		Name: labName,
		Topology: topology{
			Nodes: make(map[string]node),
			Links: make([]edge, 0),
			Kinds: map[string]kindConfig{
				"linux": {NetworkMode: "none"},
			},
		},
	}

	// indexar interfaces por ID
	ifaceLookup := make(map[string]ifaceInfo)
	for _, n := range nodesRaw {
		for _, iface := range n.Interfaces {
			ifaceLookup[iface.Id] = ifaceInfo{
				NodeName:  n.Name,
				IfaceName: iface.Name,
				Kind:      n.Kind,
			}
		}
	}

	// detectar links y conexiones a briges
	bridgeConnectionMap := make(map[string]string)

	for _, linkRaw := range edgesRaw {
		src, srcOk := ifaceLookup[linkRaw.Source]
		tgt, tgtOk := ifaceLookup[linkRaw.Target]

		if srcOk && tgtOk {
			ep1 := formatEndpoint(src)
			ep2 := formatEndpoint(tgt)

			lab.Topology.Links = append(lab.Topology.Links, edge{
				Endpoints: []string{ep1, ep2},
			})

			// registrar conexion de linux a bridge
			if src.Kind == "linux" && tgt.Kind == "bridge" {
				key := fmt.Sprintf("%s:%s", src.NodeName, src.IfaceName)
				bridgeConnectionMap[key] = tgt.NodeName
			}
			// registrar conexion de bridge a linux
			if tgt.Kind == "linux" && src.Kind == "bridge" {
				key := fmt.Sprintf("%s:%s", tgt.NodeName, tgt.IfaceName)
				bridgeConnectionMap[key] = src.NodeName
			}
		}
	}

	// configuracion de nodos
	for _, n := range nodesRaw {
		// copiar defaults
		nodeBinds := make([]string, len(defaults.Binds))
		copy(nodeBinds, defaults.Binds)

		nodeEnv := make(map[string]string)
		for k, v := range defaults.Env {
			nodeEnv[k] = v
		}

		finalNode := node{
			Kind:  n.Kind,
			Image: n.Image,
			Binds: nodeBinds,
			Env:   nodeEnv,
			Exec:  []string{},
		}

		if n.Kind == "linux" {
			finalNode.Exec = append(finalNode.Exec, fmt.Sprintf("hostname %s", n.Name))

			for _, iface := range n.Interfaces {
				key := fmt.Sprintf("%s:%s", n.Name, iface.Name)
				connectedBridgeName, isConnectedToBridge := bridgeConnectionMap[key]

				if isConnectedToBridge {
					// renombrar interfaz al nombre del bridge
					finalNode.Exec = append(finalNode.Exec,
						fmt.Sprintf("ip link set %s down", iface.Name),
						fmt.Sprintf("ip link set %s name %s", iface.Name, connectedBridgeName),
						fmt.Sprintf("ip link set %s up", connectedBridgeName),
					)

					if iface.Ip != "" {
						finalNode.Exec = append(finalNode.Exec,
							fmt.Sprintf("ip addr add %s dev %s", iface.Ip, connectedBridgeName),
						)
					}
				} else {
					// configuracion estandar
					finalNode.Exec = append(finalNode.Exec,
						fmt.Sprintf("ip link set %s up", iface.Name),
					)

					if iface.Ip != "" {
						finalNode.Exec = append(finalNode.Exec,
							fmt.Sprintf("ip addr add %s dev %s", iface.Ip, iface.Name),
						)
					}
				}
			}
		}

		lab.Topology.Nodes[n.Name] = finalNode
	}

	return yaml.Marshal(lab)
}