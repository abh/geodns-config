package main

import (
	"encoding/json"
	"log"
	"sort"
)

type zoneData map[string]*zoneLabel

type zoneJson struct {
	Data     zoneData `json:"data"`
	Ttl      int      `json:"ttl"`
	MaxHosts int      `json:"max_hosts"`
}

type jsonAddresses []interface{}

type zoneLabel struct {
	Ns    map[string]string `json:"ns,omitempty"`
	Cname string            `json:"cname,omitempty"`
	Alias string            `json:"alias,omitempty"`
	A     jsonAddresses     `json:"a,omitempty"`
	Aaaa  jsonAddresses     `json:"aaaa,omitempty"`
}

func (a jsonAddresses) Less(i, j int) bool {
	// Really this should sort on thee IP address bytes, but this is good enough
	// as we just need something to make them be in a consistent order
	return a[i].([]interface{})[0].(string) < a[j].([]interface{})[0].(string)
}
func (s jsonAddresses) Len() int      { return len(s) }
func (s jsonAddresses) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

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
	b, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (jd *zoneData) sortRecords() {
	for _, v := range *jd {
		sort.Sort(v.A)
	}
}

func (z *Zone) BuildZone() (*zoneJson, error) {
	// log.Println("BuildZone", spew.Sdump(z))

	js := zoneJson{Data: zoneData{}}

	js.MaxHosts = z.Options.MaxHosts
	js.Ttl = z.Options.Ttl

	js.Data[""] = new(zoneLabel)
	js.Data[""].Ns = map[string]string{}
	for _, ns := range z.Ns {
		js.Data[""].Ns[ns] = ""
	}

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

		if d, ok := js.Data[labelData.Name]; ok {
			if len(d.A) == 0 {
				log.Println("No global A records for", labelData.Name)
			}
		} else {
			log.Println("No global A records for", labelData.Name)
		}
	}

	js.Data.sortRecords()

	return &js, nil
}
