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
	Name string
	IP   net.IP
}

// Label has a name (hostname) and either a group name or a map of nodes
type Label struct {
	Name       string
	LabelNodes map[string]labelNode
	GroupName  string
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
		label = &Label{Name: name, LabelNodes: make(map[string]labelNode)}
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
		label = &Label{Name: name, LabelNodes: make(map[string]labelNode)}
		ls.labels[name] = label
	}

	label.LabelNodes[node.Name] = node

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

				var ip net.IP

				if len(labelTarget.(string)) > 0 {
					ip = net.ParseIP(labelTarget.(string))

					if ip == nil {
						return fmt.Errorf("Invalid IP address for '%s'/'%s': %s", name, labelName, labelTarget)
					}
				}
				node := labelNode{Name: labelName, IP: ip}
				newLabels.SetNode(name, node)
			}
		}

		ls.mutex.Lock()
		defer ls.mutex.Unlock()
		ls.labels = newLabels.labels

		return nil
	})
}
