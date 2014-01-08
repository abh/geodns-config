package dnsconfig

import (
	. "launchpad.net/gocheck"
)

type DnsSuite struct {
	zone *Zone
}

var _ = Suite(&DnsSuite{})

func (s *DnsSuite) SetUpSuite(c *C) {
	z := new(Zone)
	z.Name = "example.com"
	z.Options.Ttl = 25
	z.Labels.LoadFile("testdata/labels.json")
	z.GeoMap.LoadFile("testdata/geomap.json")
	z.Nodes.LoadFile("testdata/nodes.json")

	s.zone = z
}

func (s *DnsSuite) TestDnsLoad(c *C) {

	z := s.zone

	zd, err := z.BuildZone()
	c.Assert(err, IsNil)

	c.Check(zd.Ttl, Equals, 25)

	t1, ok := zd.Data["zone2.example"]
	c.Assert(ok, Equals, true)
	t2 := t1.A[0].([]interface{})
	// IP override, default weight
	c.Check(t2, DeepEquals, []interface{}{"10.1.1.10", 100})

	t1, ok = zd.Data["zone2.example.europe"]
	c.Assert(ok, Equals, true)
	t2 = t1.A[0].([]interface{})
	// use default IP from nodes.json and weight override
	c.Check(t2, DeepEquals, []interface{}{"10.0.5.1", 1000})

	t1, ok = zd.Data["zone4"]
	c.Assert(ok, Equals, true)

	t1, ok = zd.Data["zone4"]
	c.Assert(ok, Equals, true)

	t1, ok = zd.Data["zone4.us"]
	// edge01.sea is inactive and edge01.jfk disabled for this label
	c.Assert(ok, Equals, false)

	t1, ok = zd.Data["any-only"]
	c.Assert(ok, Equals, true)
	c.Check(t1.A[0].([]interface{}), DeepEquals, []interface{}{"10.0.0.1", 100})
	c.Check(t1.A[1].([]interface{}), DeepEquals, []interface{}{"10.0.0.4", 100})

	t1, ok = zd.Data["any-only.north-america"]
	c.Assert(ok, Equals, true)
	c.Check(t1.A[0].([]interface{}), DeepEquals, []interface{}{"10.0.0.1", 10})
	c.Check(t1.A[1].([]interface{}), DeepEquals, []interface{}{"10.0.0.4", 10})

	t1, ok = zd.Data["any-alias"]
	c.Assert(ok, Equals, true)
	c.Check(t1.Alias, Equals, "any-only")

	js, err := zd.JSON()
	c.Check(err, IsNil)
	c.Check(len(js) > 0, Equals, true)
}

func (s *DnsSuite) TestCname(c *C) {

	z := s.zone

	zd, err := z.BuildZone()
	c.Assert(err, IsNil)

	t1, ok := zd.Data["zone3.example.dk"]
	c.Assert(ok, Equals, true)
	c.Check(t1.Cname[0], DeepEquals, []interface{}{"one-override.example.com", 100})

	t1, ok = zd.Data["zone3.example.se"]
	c.Assert(ok, Equals, true)
	c.Check(t1.Cname[0], DeepEquals, []interface{}{"two.example.com", 100})

	t1, ok = zd.Data["zone3.example.no"]
	c.Assert(ok, Equals, true)
	c.Check(t1.Cname[0], DeepEquals, []interface{}{"one-override.example.com", 2})
	c.Check(t1.Cname[1], DeepEquals, []interface{}{"two.example.com", 1})

}

func (s *DnsSuite) TestDnsSort(c *C) {
	zd := zoneData{}
	l := new(zoneLabel)
	zd["test"] = l

	l.A = make([]interface{}, 4)
	l.A[0] = []interface{}{"20.2.1.4", 200}
	l.A[1] = []interface{}{"20.50.1.4", 300}
	l.A[2] = []interface{}{"1.2.3.4", 190}
	l.A[3] = []interface{}{"10.2.3.4", 150}

	zd.sortRecords()

	c.Check(l.A, DeepEquals, jsonAddresses{[]interface{}{"1.2.3.4", 190}, []interface{}{"10.2.3.4", 150}, []interface{}{"20.2.1.4", 200}, []interface{}{"20.50.1.4", 300}})
}
