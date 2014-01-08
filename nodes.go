package dnsconfig

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type nodesMap map[string]*Node

type Nodes struct {
	nodes nodesMap
	mutex sync.Mutex
}

type Node struct {
	Name   string
	Ip     net.IP
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

	for i, node := range ns.nodes {
		log.Println(i, node)
		r = append(r, node)
	}
	log.Println("Rs", r)
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

		var nodes = nodesMap{}
		for name, v := range objmap {
			data := v.(map[string]interface{})
			// log.Println("name, data", name, data)

			active, err := toBool(data["active"])
			if err != nil {
				return err
			}

			var cname string
			var ip net.IP

			if cnameIf, ok := data["cname"]; ok {
				cname = cnameIf.(string)
			} else {

				ipStr := data["ip"].(string)

				ip = net.ParseIP(ipStr)
				if ip == nil {
					return fmt.Errorf("Invalid IP address %s for node '%s'", ipStr, name)
				}
			}

			node := &Node{Cname: cname, Ip: ip, Active: active}

			nodes[name] = node
			// log.Printf("%#v\n", node)

		}

		ns.nodes = nodes

		return nil
	})
}
