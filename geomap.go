package main

import (
	"fmt"
	"strings"
	"sync"
)

type geoTargetMap map[string][]*GeoTarget

type GeoMap struct {
	geomap geoTargetMap
	mutex  sync.Mutex
}

type GeoTarget struct {
	target string
	weight int
}

func NewGeoMap() GeoMap {
	gm := GeoMap{}
	gm.Clear()
	return gm
}

func (gm *GeoMap) Clear() {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	gm.geomap = make(geoTargetMap)
}

func (gm *GeoMap) GetNodeGeos(node string) []*GeoTarget {

	for k, geos := range gm.geomap {
		if matchWildcard(k, node) {
			return geos
		}
	}

	return nil
}

func (gm *GeoMap) LoadFile(fileName string) error {

	objmap := make(map[string]interface{})

	return jsonLoader(fileName, objmap, func() error {
		gm.mutex.Lock()
		defer gm.mutex.Unlock()

		geomap := geoTargetMap{}

		for name, v := range objmap {
			// log.Printf("%s: %#v\n", name, v)

			if _, ok := geomap[name]; !ok {
				geomap[name] = make([]*GeoTarget, 0)
			}

			for _, g := range v.([]interface{}) {
				gSplit := strings.Split(g.(string), "=")

				// log.Printf("gSplit: %#v\n", gSplit)

				weight := 100

				if len(gSplit) > 1 {
					var err error
					weight, err = toInt(gSplit[1])
					if err != nil {
						return fmt.Errorf("Bad weight '%s' for geo '%s'/'%s': %s\n", gSplit[1], name, gSplit[0], err)
					}
				}

				geo := GeoTarget{target: gSplit[0], weight: weight}
				geomap[name] = append(geomap[name], &geo)

			}

			// log.Printf("%s: %#v", name, geomap)

		}

		gm.geomap = geomap

		return nil
	})
}
