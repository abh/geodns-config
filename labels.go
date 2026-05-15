package dnsconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/netip"
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
	IP     netip.Addr
	Cname  string
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
			return &labelNode{Name: name, Active: node.Active, IP: node.IP, Cname: node.Cname}
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

// labelEntry accepts either an IP string shorthand ("10.1.2.3") or
// an object {ip, cname, active}.
type labelEntry struct {
	IP     string    `json:"ip"`
	Cname  string    `json:"cname"`
	Active *flexBool `json:"active"`
}

func (e *labelEntry) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '"' {
		return json.Unmarshal(data, &e.IP)
	}
	type alias labelEntry
	return json.Unmarshal(data, (*alias)(e))
}

// labelDoc decodes one label: an optional "group" alias plus zero or
// more node-name entries.
type labelDoc struct {
	Group string
	Nodes map[string]labelEntry
}

func (d *labelDoc) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Nodes = make(map[string]labelEntry, len(raw))
	for k, v := range raw {
		if k == "group" {
			if err := json.Unmarshal(v, &d.Group); err != nil {
				return fmt.Errorf("invalid group: %w", err)
			}
			continue
		}
		var entry labelEntry
		if err := json.Unmarshal(v, &entry); err != nil {
			return fmt.Errorf("invalid entry '%s': %w", k, err)
		}
		d.Nodes[k] = entry
	}
	return nil
}

// LoadFile loads a labels.json file, replacing any existing labels.
func (ls *Labels) LoadFile(fileName string) error {
	var raw map[string]labelDoc
	if err := loadJSONFile(fileName, &raw); err != nil {
		return err
	}

	newLabels := NewLabels()
	for labelName, doc := range raw {
		if doc.Group != "" {
			newLabels.SetGroup(labelName, doc.Group)
		}
		for nodeName, entry := range doc.Nodes {
			node := labelNode{Name: nodeName, Active: true, Cname: entry.Cname}
			if entry.Active != nil {
				node.Active = bool(*entry.Active)
			}
			if entry.IP != "" {
				ip, err := netip.ParseAddr(entry.IP)
				if err != nil {
					return fmt.Errorf("invalid IP address for '%s'/'%s': %s: %w", labelName, nodeName, entry.IP, err)
				}
				node.IP = ip
			}
			newLabels.SetNode(labelName, node)
		}
	}

	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.labels = newLabels.labels
	return nil
}
