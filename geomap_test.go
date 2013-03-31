package dnsconfig

import (
	. "launchpad.net/gocheck"
)

type GeoMapSuite struct {
	GeoMap GeoMap
}

var _ = Suite(&GeoMapSuite{})

func (s *GeoMapSuite) SetUpSuite(c *C) {
	s.GeoMap = NewGeoMap()
}

func (s *GeoMapSuite) TestGeoMap(c *C) {
	s.GeoMap.Clear()
}

func (s *GeoMapSuite) TestGeoLoad(c *C) {
	s.GeoMap.Clear()
	err := s.GeoMap.LoadFile("testdata/geomap.json")
	c.Assert(err, IsNil)

	// results are sorted appropriately
	c.Assert(s.GeoMap.geomap["*.ams*"][0].target, Equals, "fr")
	c.Assert(s.GeoMap.geomap["*.ams*"][1].target, Equals, "nl")
	c.Assert(s.GeoMap.geomap["*.ams*"][2].target, Equals, "europe")

	// "@" gets sorted first
	c.Assert(s.GeoMap.GetNodeGeos("test.sea")[0].target, Equals, "@")

	c.Assert(s.GeoMap.geomap["*.lhr"][1].weight, Equals, 1000)

	// make sure we get the more specific entry
	c.Assert(s.GeoMap.GetNodeGeos("flex04.ams04")[0].weight, Equals, 1)
	c.Assert(s.GeoMap.GetNodeGeos("flex04.ams04")[0].target, Equals, "europe")

	c.Assert(s.GeoMap.GetNodeGeos("x123.lhr")[1].weight, Equals, 1000)
	c.Assert(s.GeoMap.GetNodeGeos("x123.lhr")[1].target, Equals, "europe")

	c.Assert(s.GeoMap.GetNodeGeos("x123.lhr")[0].weight, Equals, 100)

}
