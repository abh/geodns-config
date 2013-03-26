package main

import (
	"encoding/json"
	"log"
)

type zoneJson struct {
	Data     zoneData `json:"data"`
	Ttl      int      `json:"ttl"`
	MaxHosts int      `json:"max_hosts"`
}

type zoneLabel struct {
	Ns    map[string]string `json:"ns,omitempty"`
	Cname string            `json:"cname,omitempty"`
	Alias string            `json:"alias,omitempty"`
	A     []interface{}     `json:"a,omitempty"`
	Aaaa  []interface{}     `json:"aaaa,omitempty"`
}

type zoneData map[string]*zoneLabel

func (z *Zone) BuildJSON() (string, error) {
	zd, err := z.BuildZone()
	if err != nil {
		return "", err
	}
	js, err := zd.JSON()
	if err != nil {
		return "", err
	}
	return js, nil
}

func (js *zoneJson) JSON() (string, error) {
	b, err := json.Marshal(js)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (z *Zone) BuildZone() (*zoneJson, error) {
	// log.Println("BuildZone", spew.Sdump(z))

	js := zoneJson{Data: zoneData{}}

	js.MaxHosts = z.Options.MaxHosts
	js.Ttl = z.Options.Ttl

	js.Data[""] = new(zoneLabel)
	js.Data[""].Ns = map[string]string{"a.ntpns.org": "", "b.ntpns.org": ""}

	for _, labelData := range z.Labels.All() {
		if len(labelData.GroupName) > 0 {
			log.Println("Groups not implemented yet, skipping ", labelData.Name)
			continue
		}
		for _, labelNode := range labelData.LabelNodes {
			node := z.Nodes.Get(labelNode.Name)
			if node == nil {
				log.Printf("Node '%s' not configured in master nodes config\n", labelNode.Name)
				continue
			}
			if node.Active == false {
				log.Printf("Node '%s' is inactive for label '%s'\n", labelNode.Name, labelData.Name)
				continue
			}

			geos := z.GeoMap.GetNodeGeos(labelNode.Name)

			for _, geo := range geos {
				var geoName string
				if geo.target == "@" {
					geoName = labelData.Name
				} else {
					geoName = labelData.Name + "." + geo.target
				}
				if _, ok := js.Data[geoName]; !ok {
					js.Data[geoName] = new(zoneLabel)
				}

				ip := labelNode.Ip
				if ip == nil {
					ip = node.Ip
				}

				trg := []interface{}{ip.String(), geo.weight}
				js.Data[geoName].A = append(js.Data[geoName].A, trg)
			}

		}
	}

	return &js, nil
}
