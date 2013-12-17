package dnsconfig

import (
	"fmt"
	"net"
	"sync"
)

type labelsMap map[string]*Label

// Labels map hostnames to nodes or groups
type Labels struct {
	labels labelsMap
	mutex  sync.Mutex
}

type labelNode struct {
	Name   string
	Active bool
	IP     net.IP
}

// Label has a name (hostname) and either a group name or a map of nodes
type Label struct {
	Name       string
	labelNodes map[string]labelNode
	GroupName  string
}

// Get all nodes for a label
func (l *Label) GetNodes() []labelNode {
	lbls := make([]labelNode, 0)
	for _, lbl := range l.labelNodes {
		lbls = append(lbls, lbl)
	}
	return lbls
}

func (l *Label) GetNode(name string) *labelNode {
	lbl := l.labelNodes[name]
	if lbl.Name != "" {
		return &lbl
	}
	for nodeName, node := range l.labelNodes {
		if matchWildcard(nodeName, name) {
			return &labelNode{Name: name, Active: node.Active, IP: node.IP}
		}
	}
	return nil
}

// NewLabels return a new Labels struct
func NewLabels() Labels {
	ls := Labels{}
	ls.Clear()
	return ls
}

// All returns a slice of the labels
func (ls *Labels) All() (r []*Label) {

	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	for _, label := range ls.labels {
		r = append(r, label)
	}
	return
}

// Clear resets the labels map
func (ls *Labels) Clear() {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.labels = labelsMap{}
}

// SetGroup sets the label to be a particular group
func (ls *Labels) SetGroup(name, groupName string) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	label, ok := ls.labels[name]
	if !ok {
		label = &Label{Name: name, labelNodes: make(map[string]labelNode)}
		ls.labels[name] = label
	}
	label.GroupName = groupName
}

// SetNode adds a node to a label
func (ls *Labels) SetNode(name string, node labelNode) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	label, ok := ls.labels[name]
	if !ok {
		label = &Label{Name: name, labelNodes: make(map[string]labelNode)}
		ls.labels[name] = label
	}

	label.labelNodes[node.Name] = node

}

// Get returns a named label
func (ls *Labels) Get(name string) *Label {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	if label, ok := ls.labels[name]; ok {
		return label
	}
	return nil
}

// Count returns the number of labels
func (ls *Labels) Count() int {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	return len(ls.labels)
}

// LoadFile loads a labels.json file into the data structure. It is not currently
// cleared first.
func (ls *Labels) LoadFile(fileName string) error {

	objmap := make(objMap)

	return jsonLoader(fileName, objmap, func() error {
		var newLabels = NewLabels()

		for name, v := range objmap {
			data := v.(map[string]interface{})

			for labelName, labelTarget := range data {

				if labelName == "group" {
					newLabels.SetGroup(name, labelTarget.(string))
					continue
				}

				// fmt.Printf("labelName '%s', labelTarget '%#v'\n", labelName, labelTarget)

				node := labelNode{Name: labelName, Active: true}

				var ipStr string

				switch labelTarget.(type) {
				case string:
					ipStr = labelTarget.(string)
				case map[string]interface{}:
					v := labelTarget.(map[string]interface{})
					if ipV, ok := v["ip"]; ok {
						ipStr = ipV.(string)
					}
					if activeV, ok := v["active"]; ok {
						active, err := toBool(activeV)
						node.Active = active
						if err != nil {
							fmt.Errorf("Invalid active flag for '%s'/'%s': %s", name, labelName, active)
						}
					}

				default:
					return fmt.Errorf("Invalid value type for '%s'/%s': %T (%#v)", name, labelName, labelTarget, labelTarget)
				}

				if len(ipStr) > 0 {
					ip := net.ParseIP(ipStr)
					if ip == nil {
						return fmt.Errorf("Invalid IP address for '%s'/'%s': %s", name, labelName, ipStr)
					}
					node.IP = ip
				}

				newLabels.SetNode(name, node)
			}
		}

		ls.mutex.Lock()
		defer ls.mutex.Unlock()
		ls.labels = newLabels.labels

		return nil
	})
}
