package main

import (
	. "launchpad.net/gocheck"
)

type DnsSuite struct {
}

var _ = Suite(&DnsSuite{})

func (s *DnsSuite) SetUpSuite(c *C) {
}

func (s *LabelsSuite) TestDnsLoad(c *C) {
	z := &Zone{Name: "example.com"}
	z.Labels.LoadFile("testdata/labels.json")
	z.GeoMap.LoadFile("testdata/geomap.json")
	z.Nodes.LoadFile("testdata/nodes.json")

	zd, err := z.BuildZone()
	c.Assert(err, IsNil)

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

	t1, ok = zd.Data["zone4.us"]
	c.Assert(ok, Equals, false) // edge01.sea is inactive

	// log.Println("T1", spew.Sdump(t1))

	// c.Assert(t1, Equals, "10.1.1.10")

	js, err := zd.JSON()
	c.Check(err, IsNil)
	c.Check(len(js) > 0, Equals, true)
}
