package main

import (
	"camlistore.org/pkg/errorutil"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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

	fh, err := os.Open(fileName)
	if err != nil {
		log.Println("Could not read ", fileName, ": ", err)
		return err
	}

	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	var objmap map[string]interface{}

	decoder := json.NewDecoder(fh)
	if err = decoder.Decode(&objmap); err != nil {
		extra := ""
		if serr, ok := err.(*json.SyntaxError); ok {
			if _, serr := fh.Seek(0, os.SEEK_SET); serr != nil {
				log.Fatalf("seek error: %v", serr)
			}
			line, col, highlight := errorutil.HighlightBytePosition(fh, serr.Offset)
			extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
				line, col, serr.Offset, highlight)
		}
		return fmt.Errorf("error parsing JSON object in config file %s%s\n%v",
			fh.Name(), extra, err)
	}

	var nodes = nodesMap{}
	for name, v := range objmap {
		data := v.(map[string]interface{})
		log.Println("name, data", name, data)

		active, err := toBool(data["active"])
		if err != nil {
			return err
		}

		ip := net.ParseIP(data["ip"].(string))
		if ip == nil {
			return fmt.Errorf("Invalid IP address %s", data["ip"].(string))
		}

		node := &Node{Ip: ip, Active: active}

		nodes[name] = node
		log.Printf("%#v\n", node)

	}

	ns.nodes = nodes

	return nil
}

func toInt(i interface{}) (int, error) {
	switch i.(type) {
	case string:
		return strconv.Atoi(i.(string))
	case float64:
		return int(i.(float64)), nil
	}
	return 0, fmt.Errorf("Unknown type %T", i)
}

func toBool(i interface{}) (bool, error) {
	n, err := toInt(i)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
