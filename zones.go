package dnsconfig

import (
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

func (zs *Zones) LoadZonesConfig(fileName string) error {

	objmap := objMap{}

	return jsonLoader(fileName, objmap, func() error {

		zs.mutex.Lock()
		defer zs.mutex.Unlock()

		if zs.zones == nil {
			zs.zones = map[string]*Zone{}
		}

		for zoneName, zoneData := range objmap {
			var zone *Zone

			if zone2, ok := zs.zones[zoneName]; ok {
				zone = zone2
			} else {
				zone = new(Zone)
				zs.zones[zoneName] = zone
				zone.Name = zoneName
				zone.Options.Ttl = 300
				zone.Options.Contact = "hostmaster"
			}

			zoneOptions := zoneData.(map[string]interface{})

			for key, v := range zoneOptions {
				switch key {
				case "ttl":
					i, err := toInt(v)
					if err != nil {
						return fmt.Errorf("Invalid integer '%s' in '%s' option: %s", v, key, err)
					}
					zone.Options.Ttl = i
				case "serial":
					i, err := toInt(v)
					if err != nil {
						return fmt.Errorf("Invalid integer '%s' in '%s' option: %s", v, key, err)
					}
					zone.Options.Serial = i
				case "contact":
					zone.Options.Contact = v.(string)
				case "max_hosts":
					i, err := toInt(v)
					if err != nil {
						return fmt.Errorf("Invalid integer '%s' in '%s' option: %s", v, key, err)
					}
					zone.Options.MaxHosts = i
				case "ns":
					switch v.(type) {
					case []interface{}:
						nsList := v.([]interface{})
						ns := make([]string, len(nsList))
						for i, v := range nsList {
							ns[i] = v.(string)
						}
						zone.Ns = ns
					default:
						return fmt.Errorf("Bad ns parameter for '%s'\n", zoneName)
					}

				case "labels":
					zone.LabelsFile = absPath(fileName, v.(string))
				case "nodes":
					zone.NodesFile = absPath(fileName, v.(string))
				case "geomap":
					zone.GeoMapFile = absPath(fileName, v.(string))

				case "targeting":
					zone.Options.Targeting = v.(string)

				case "logging":
					m := v.(map[string]interface{})
					if o, ok := m["stathat"]; ok {
						zone.Logging.StatHat = o.(bool)
					}
					if o, ok := m["stathat_api"]; ok {
						zone.Logging.StatHatAPI = o.(string)
					}

				default:
					log.Printf("Unknown option '%s' for zone '%s'\n", key, zoneName)
				}
			}
		}

		return nil
	})
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
