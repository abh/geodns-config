package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type labelsMap map[string]Label

type Labels struct {
	labels labelsMap
	mutex  sync.Mutex
}

type labelNode struct {
	Name string
	Ip   net.IP
}

type Label struct {
	Name       string
	LabelNodes map[string]labelNode
	GroupName  string
}

func NewLabels() Labels {
	ls := Labels{}
	ls.Clear()
	return ls
}

func (ls *Labels) All() (r []*Label) {

	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	// for i, node := range ls.labels {
	// log.Println(i, node)
	// r = append(r, node)
	// }
	// log.Println("Rs", r)
	return
}

func (ls *Labels) Clear() {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.labels = labelsMap{}
}

func (ls *Labels) SetGroup(name, groupName string) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	label, ok := ls.labels[name]
	if !ok {
		label = Label{Name: name, LabelNodes: make(map[string]labelNode)}
		ls.labels[name] = label
	}
	label.GroupName = groupName
}

func (ls *Labels) SetNode(name string, node labelNode) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	label, ok := ls.labels[name]
	if !ok {
		label = Label{Name: name, LabelNodes: make(map[string]labelNode)}
		ls.labels[name] = label
	}

	label.LabelNodes[node.Name] = node

}

func (ls *Labels) Get(name string) *Label {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	if label, ok := ls.labels[name]; ok {
		return &label
	}
	return nil
}

func (ls *Labels) Count() int {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	return len(ls.labels)
}

func (ls *Labels) LoadFile(fileName string) error {

	objmap := make(objMap)

	return jsonLoader(fileName, objmap, func() error {

		log.Println("Loading labels from", fileName)

		var newLabels = NewLabels()

		for name, v := range objmap {
			data := v.(map[string]interface{})

			for labelName, labelTarget := range data {

				if labelName == "group" {
					newLabels.SetGroup(name, labelTarget.(string))
					continue
				}

				ip := net.ParseIP(labelTarget.(string))
				if ip == nil {
					return fmt.Errorf("Invalid IP address for '%s'/'%s': %s", name, labelName, labelTarget)
				}

				node := labelNode{Name: labelName, Ip: ip}
				newLabels.SetNode(name, node)
			}
		}

		ls.mutex.Lock()
		defer ls.mutex.Unlock()
		ls.labels = newLabels.labels

		return nil
	})
}
