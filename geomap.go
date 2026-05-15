package dnsconfig

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type geoTargetList []*GeoTarget

type (
	geoTargetMap map[string]geoTargetList
	GeoMap       struct {
		geomap geoTargetMap
		mutex  sync.Mutex
	}
)

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

	// Map iteration is unordered, so collect candidates and sort them
	// for deterministic results when multiple wildcards match.
	// Longer patterns (more literal characters) win; alphabetical as tiebreak.
	var matches []string
	for k := range gm.geomap {
		if matchWildcard(k, node) {
			matches = append(matches, k)
		}
	}
	if len(matches) > 0 {
		sort.Slice(matches, func(i, j int) bool {
			if len(matches[i]) != len(matches[j]) {
				return len(matches[i]) > len(matches[j])
			}
			return matches[i] < matches[j]
		})
		return gm.geomap[matches[0]]
	}

	if geo, ok := gm.geomap["default"]; ok {
		return geo
	}

	return nil
}

func (gm *GeoMap) LoadFile(fileName string) error {
	var raw map[string][]string
	if err := loadJSONFile(fileName, &raw); err != nil {
		return err
	}

	geomap := geoTargetMap{}
	for name, targets := range raw {
		list := make(geoTargetList, 0, len(targets))
		for _, g := range targets {
			target, weightStr, hasWeight := strings.Cut(g, "=")
			weight := 100
			if hasWeight {
				w, err := strconv.Atoi(weightStr)
				if err != nil {
					return fmt.Errorf("bad weight '%s' for geo '%s'/'%s': %w", weightStr, name, target, err)
				}
				weight = w
			}
			list = append(list, &GeoTarget{target: target, weight: weight})
		}
		list.Sort()
		geomap[name] = list
	}

	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	gm.geomap = geomap
	return nil
}
