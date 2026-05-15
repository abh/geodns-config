package dnsconfig

import (
	"fmt"
	"net/netip"
	"sync"
)

type nodesMap map[string]*Node

type Nodes struct {
	nodes nodesMap
	mutex sync.Mutex
}

type Node struct {
	Name   string
	IP     netip.Addr
	Cname  string
	Active bool
}

func NewNodes() Nodes {
	ns := Nodes{}
	ns.Clear()
	return ns
}

func (ns *Nodes) All() (r []*Node) {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	for _, node := range ns.nodes {
		r = append(r, node)
	}
	return
}

func (ns *Nodes) Clear() {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	ns.nodes = nodesMap{}
}

func (ns *Nodes) Set(name string, node Node) {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	node.Name = name
	ns.nodes[name] = &node
}

func (ns *Nodes) Get(name string) *Node {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	if node, ok := ns.nodes[name]; ok {
		return node
	}
	return nil
}

func (ns *Nodes) Count() int {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	return len(ns.nodes)
}

type nodeJSON struct {
	IP     string   `json:"ip"`
	Cname  string   `json:"cname"`
	Active flexBool `json:"active"`
}

func (ns *Nodes) LoadFile(fileName string) error {
	var raw map[string]nodeJSON
	if err := loadJSONFile(fileName, &raw); err != nil {
		return err
	}

	nodes := nodesMap{}
	for name, data := range raw {
		node := &Node{Name: name, Cname: data.Cname, Active: bool(data.Active)}
		if data.Cname == "" {
			ip, err := netip.ParseAddr(data.IP)
			if err != nil {
				return fmt.Errorf("invalid IP address %s for node '%s': %w", data.IP, name, err)
			}
			node.IP = ip
		}
		nodes[name] = node
	}

	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	ns.nodes = nodes
	return nil
}
