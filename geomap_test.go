package main

import (
	. "launchpad.net/gocheck"
	"log"
)

type GeoMapSuite struct {
	GeoMap GeoMap
}

var _ = Suite(&GeoMapSuite{})

func (s *GeoMapSuite) SetUpSuite(c *C) {
	s.GeoMap = NewGeoMap()
	log.Println("setup geomap", s)
}

func (s *GeoMapSuite) TestGeoMap(c *C) {
	s.GeoMap.Clear()
}

func (s *GeoMapSuite) TestLoad(c *C) {
	s.GeoMap.Clear()
	err := s.GeoMap.LoadFile("testdata/geomap.json")
	c.Assert(err, IsNil)

	c.Assert(s.GeoMap.geomap["*.ams"][0].target, Equals, "europe")
	c.Assert(s.GeoMap.geomap["*.lhr"][0].weight, Equals, 1000)
}
