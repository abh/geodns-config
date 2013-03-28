package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type geoTargetList []*GeoTarget

type geoTargetMap map[string]geoTargetList
type GeoMap struct {
	geomap geoTargetMap
	mutex  sync.Mutex
}

type GeoTarget struct {
	target string
	weight int
}

func (a geoTargetList) Less(i, j int) bool {
	iLen := len(a[i].target)
	jLen := len(a[j].target)
	if iLen == jLen {
		return a[i].target < a[j].target
	}
	return iLen < jLen
}
func (s geoTargetList) Len() int      { return len(s) }
func (s geoTargetList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s geoTargetList) Sort()         { sort.Sort(s) }

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

	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if geo, ok := gm.geomap[node]; ok {
		return geo
	}

	for k, geos := range gm.geomap {
		if matchWildcard(k, node) {
			return geos
		}
	}

	if geo, ok := gm.geomap["default"]; ok {
		return geo
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

			geomap[name].Sort()

			// log.Printf("%s: %#v", name, geomap)

		}

		gm.geomap = geomap

		return nil
	})
}
