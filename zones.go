package dnsconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sync"
)

type Zones struct {
	mutex sync.Mutex
	zones map[string]*Zone
}

type Zone struct {
	Name       string
	Options    ZoneOptions
	Logging    ZoneLogging
	Ns         []string
	Labels     Labels
	Nodes      Nodes
	GeoMap     GeoMap
	LabelsFile string
	NodesFile  string
	GeoMapFile string
	Verbose    bool
}

type ZoneOptions struct {
	Serial    int
	Ttl       int
	MaxHosts  int
	Targeting string
	Contact   string
}

func (z *Zone) LoadConfig() error {
	err := z.Nodes.LoadFile(z.NodesFile)
	if err != nil {
		return err
	}
	err = z.GeoMap.LoadFile(z.GeoMapFile)
	if err != nil {
		return err
	}
	err = z.Labels.LoadFile(z.LabelsFile)
	if err != nil {
		return err
	}
	return nil
}

func (zs *Zones) All() (r []*Zone) {
	zs.mutex.Lock()
	defer zs.mutex.Unlock()

	for _, zone := range zs.zones {
		r = append(r, zone)
	}
	return
}

type zoneJSON struct {
	TTL       flexInt     `json:"ttl"`
	Serial    flexInt     `json:"serial"`
	Contact   string      `json:"contact"`
	MaxHosts  flexInt     `json:"max_hosts"`
	NS        []string    `json:"ns"`
	Labels    string      `json:"labels"`
	Nodes     string      `json:"nodes"`
	GeoMap    string      `json:"geomap"`
	Targeting string      `json:"targeting"`
	Logging   loggingJSON `json:"logging"`
}

type loggingJSON struct {
	StatHat    bool   `json:"stathat"`
	StatHatAPI string `json:"stathat_api"`
}

var validZoneKeys = map[string]bool{
	"ttl": true, "serial": true, "contact": true, "max_hosts": true,
	"ns": true, "labels": true, "nodes": true, "geomap": true,
	"targeting": true, "logging": true,
}

func (zs *Zones) LoadZonesConfig(fileName string) error {
	var raw map[string]json.RawMessage
	if err := loadJSONFile(fileName, &raw); err != nil {
		return err
	}

	zs.mutex.Lock()
	defer zs.mutex.Unlock()
	if zs.zones == nil {
		zs.zones = map[string]*Zone{}
	}

	for zoneName, rawZone := range raw {
		var keys map[string]json.RawMessage
		if err := json.Unmarshal(rawZone, &keys); err != nil {
			return fmt.Errorf("invalid zone '%s': %w", zoneName, err)
		}

		zone, ok := zs.zones[zoneName]
		if !ok {
			zone = &Zone{Name: zoneName, Options: ZoneOptions{Ttl: 300, Contact: "hostmaster"}}
			zs.zones[zoneName] = zone
		}

		// Pre-seed with current values; absent JSON keys leave them untouched.
		data := zoneJSON{
			TTL:       flexInt(zone.Options.Ttl),
			Serial:    flexInt(zone.Options.Serial),
			Contact:   zone.Options.Contact,
			MaxHosts:  flexInt(zone.Options.MaxHosts),
			NS:        zone.Ns,
			Targeting: zone.Options.Targeting,
			Logging:   loggingJSON(zone.Logging),
		}
		if err := json.Unmarshal(rawZone, &data); err != nil {
			return fmt.Errorf("invalid zone '%s': %w", zoneName, err)
		}

		for k := range keys {
			if !validZoneKeys[k] {
				log.Printf("Unknown option '%s' for zone '%s'\n", k, zoneName)
			}
		}

		zone.Options.Ttl = int(data.TTL)
		zone.Options.Serial = int(data.Serial)
		zone.Options.Contact = data.Contact
		zone.Options.MaxHosts = int(data.MaxHosts)
		zone.Options.Targeting = data.Targeting
		zone.Logging = ZoneLogging(data.Logging)
		zone.Ns = data.NS
		if data.Labels != "" {
			zone.LabelsFile = absPath(fileName, data.Labels)
		}
		if data.Nodes != "" {
			zone.NodesFile = absPath(fileName, data.Nodes)
		}
		if data.GeoMap != "" {
			zone.GeoMapFile = absPath(fileName, data.GeoMap)
		}
	}
	return nil
}

func absPath(baseConfig, fileName string) string {
	abs, err := filepath.Abs(baseConfig)
	if err != nil {
		log.Println("Could not determine absolute path for", baseConfig)
		return fileName
	}
	dir := filepath.Dir(abs)
	return filepath.Join(dir, fileName)
}
