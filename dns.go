package dnsconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type ZoneLogging struct {
	StatHat    bool   `json:"stathat"`
	StatHatAPI string `json:"stathat_api"`
}

type zoneData map[string]*zoneLabel

type zoneJson struct {
	Data      zoneData    `json:"data"`
	Ttl       int         `json:"ttl"`
	MaxHosts  int         `json:"max_hosts,omitempty"`
	Logging   ZoneLogging `json:"logging,omitempty"`
	Targeting string      `json:"targeting,omitempty"`
}

type jsonAddresses []interface{}

type zoneLabel struct {
	Ns    map[string]string `json:"ns,omitempty"`
	Cname jsonAddresses     `json:"cname,omitempty"`
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
	js.Logging = z.Logging
	js.Targeting = z.Options.Targeting

	js.Data[""] = new(zoneLabel)
	js.Data[""].Ns = map[string]string{}
	for _, ns := range z.Ns {
		js.Data[""].Ns[ns] = ""
	}

	displayQueue := make([]string, 0)

	for _, labelData := range z.Labels.All() {
		if len(labelData.GroupName) > 0 {
			label := new(zoneLabel)
			label.Alias = labelData.GroupName
			js.Data[labelData.Name] = label
			continue
		}
		for _, labelNode := range labelData.GetNodes() {

			if labelNode.Active == false {
				log.Printf("Node '%s' is inactive in label '%s'\n", labelNode.Name, labelData.Name)
				continue
			}

			node := z.Nodes.Get(labelNode.Name)
			if node == nil {
				log.Printf("Node '%s' not configured in master nodes config\n", labelNode.Name)
				continue
			}

			geos := z.GeoMap.GetNodeGeos(labelNode.Name)

			// Don't warn if there are no targets for the inactive node
			if node.Active == false && len(geos) > 0 {
				log.Printf("Node '%s' is inactive (used in '%s')\n", labelNode.Name, labelData.Name)
				continue
			}

			for _, geo := range geos {
				var geoName string
				if geo.target == "@" {
					geoName = labelData.Name
				} else {
					if len(labelData.Name) > 0 {
						geoName = labelData.Name + "." + geo.target
					} else {
						geoName = geo.target
					}
				}
				if _, ok := js.Data[geoName]; !ok {
					js.Data[geoName] = new(zoneLabel)
				}

				cname := labelNode.Cname
				if len(cname) == 0 {
					cname = node.Cname
				}
				if len(cname) > 0 {
					trg := []interface{}{cname, geo.weight}
					js.Data[geoName].Cname = append(js.Data[geoName].Cname, trg)
				} else {
					ip := labelNode.IP
					if ip == nil {
						ip = node.Ip
					}

					trg := []interface{}{ip.String(), geo.weight}
					js.Data[geoName].A = append(js.Data[geoName].A, trg)
				}

				fn := func(slice []string, s string) []string {
					for _, e := range slice {
						if e == s {
							return slice
						}
					}
					return append(slice, s)
				}
				displayQueue = fn(displayQueue, geoName)
			}
		}
		if d, ok := js.Data[labelData.Name]; ok {
			if len(d.A) == 0 && len(d.Cname) == 0 {
				log.Println("No global data for", labelData.Name)
			}
		} else {
			log.Println("No global data for", labelData.Name)
		}
		displayQueue = append(displayQueue, "")
	}

	js.Data.sortRecords()

	if z.Verbose {
		for _, geoName := range displayQueue {
			if geoName == "" {
				fmt.Println("")
				continue
			}
			fmt.Printf("%-40s: ", geoName)

			if len(js.Data[geoName].Cname) > 0 {
				fmt.Printf("%s\n", js.Data[geoName].Cname)
			} else {

				for i, a := range js.Data[geoName].A {
					// fmt.Printf("%#v\n%s\n", a, spew.Sdump(a))
					fmt.Printf("%-15s/%-4d", a.([]interface{})[0].(string), a.([]interface{})[1].(int))
					if i == len(js.Data[geoName].A)-1 {
						fmt.Printf("\n")
					} else {
						fmt.Printf(" | ")
					}
				}
			}
		}
	}

	return &js, nil
}
