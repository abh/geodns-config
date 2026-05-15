package dnsconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"slices"
	"strings"
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

func recordString(r interface{}) string {
	return r.([]interface{})[0].(string)
}

func compareIP(a, b interface{}) int {
	aAddr, _ := netip.ParseAddr(recordString(a))
	bAddr, _ := netip.ParseAddr(recordString(b))
	return aAddr.Compare(bAddr)
}

func compareString(a, b interface{}) int {
	return strings.Compare(recordString(a), recordString(b))
}

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
		slices.SortFunc(v.A, compareIP)
		slices.SortFunc(v.Aaaa, compareIP)
		slices.SortFunc(v.Cname, compareString)
	}
}

func (z *Zone) BuildZone() (*zoneJson, error) {
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

			if !labelNode.Active {
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
			if !node.Active && len(geos) > 0 {
				log.Printf("Node '%s' is inactive (used in '%s')\n", labelNode.Name, labelData.Name)
				continue
			}

			for _, geo := range geos {
				var geoName string
				switch {
				case geo.target == "@":
					geoName = labelData.Name
				case len(labelData.Name) > 0:
					geoName = labelData.Name + "." + geo.target
				default:
					geoName = geo.target
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
					if !ip.IsValid() {
						ip = node.IP
					}

					trg := []interface{}{ip.String(), geo.weight}
					if ip.Is4() || ip.Is4In6() {
						js.Data[geoName].A = append(js.Data[geoName].A, trg)
					} else {
						js.Data[geoName].Aaaa = append(js.Data[geoName].Aaaa, trg)
					}
				}

				if z.Verbose && !slices.Contains(displayQueue, geoName) {
					displayQueue = append(displayQueue, geoName)
				}
			}
		}
		d, ok := js.Data[labelData.Name]
		if !ok || (len(d.A) == 0 && len(d.Aaaa) == 0 && len(d.Cname) == 0) {
			log.Println("No global data for", labelData.Name)
		}
		if z.Verbose {
			displayQueue = append(displayQueue, "")
		}
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
				addrs := slices.Concat(js.Data[geoName].A, js.Data[geoName].Aaaa)
				for i, a := range addrs {
					fmt.Printf("%-39s/%-4d", a.([]interface{})[0].(string), a.([]interface{})[1].(int))
					if i == len(addrs)-1 {
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
