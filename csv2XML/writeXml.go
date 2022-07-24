package csv2xml

import (
	"encoding/xml"
	"fmt"
	"os"
)

type maps struct {
	XMLName xml.Name `xml:"map"`
	Version string   `xml:"version,attr"`
	Nodes   []*Node   `xml:"node"`
}

func WriteXML(csv *Csv) {
	nodeSeen := make(map[string]*Node)

	file, err := os.Create(fmt.Sprint(csv.FileDir, "/", csv.Filename, ".mm"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file error: %v", err)
	}

	defer file.Close()

	root := &maps{Version: "freeplane 1.9.13"}
	for _, node := range csv.NodeList {
		if node.ParentId == "" {
			root.Nodes = append(root.Nodes, node)
			nodeSeen[node.Id] = node
		} else {
			parentNode := nodeSeen[node.ParentId]
			parentNode.Nodes = append(parentNode.Nodes, node)
			nodeSeen[node.Id] = node
		}
	}

	output, err := xml.MarshalIndent(root, "", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}

	file.Write(output)

}
