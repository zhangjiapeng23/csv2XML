package csv2xml

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

type Node struct {
	XMLName xml.Name `xml:"node"`
	Name string `xml:"TEXT,attr"`
	Id string `xml:"ID,attr"`
	ParentId string `xml:"-"`
	Order int `xml:"-"`
	Created string `xml:"CREATED,attr"`
	Modified string `xml:"MODIFIED,attr"`
	Nodes []*Node `xml:"node"`
}


func NewNode(name string) *Node {
	return &Node{
		Name: name,
		Id: fmt.Sprint("ID_", uuid.NewV4().String()),
		ParentId: "",
		Order: 0,
		Created: fmt.Sprint(time.Now().Unix()),
		Modified: fmt.Sprint(time.Now().Unix()),
		Nodes: make([]*Node, 0),
	}

}

func (node *Node) String() string {
	return fmt.Sprintln("Name:", node.Name, "Id:", node.Id, "ParentId:", node.ParentId)
}