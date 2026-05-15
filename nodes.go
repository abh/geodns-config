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

func (ns *Nodes) LoadFile(fileName string) error {
	objmap := make(objMap)

	return jsonLoader(fileName, objmap, func() error {
		ns.mutex.Lock()
		defer ns.mutex.Unlock()

		nodes := nodesMap{}
		for name, v := range objmap {
			data := v.(map[string]interface{})

			active, err := toBool(data["active"])
			if err != nil {
				return err
			}

			var cname string
			var ip netip.Addr

			if cnameIf, ok := data["cname"]; ok {
				cname = cnameIf.(string)
			} else {
				ipStr := data["ip"].(string)
				ip, err = netip.ParseAddr(ipStr)
				if err != nil {
					return fmt.Errorf("invalid IP address %s for node '%s': %w", ipStr, name, err)
				}
			}

			nodes[name] = &Node{Cname: cname, IP: ip, Active: active}
		}

		ns.nodes = nodes

		return nil
	})
}
