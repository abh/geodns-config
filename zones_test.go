package main

import (
	"github.com/davecgh/go-spew/spew"
	. "launchpad.net/gocheck"
)

type ZonesSuite struct {
}

var _ = Suite(&ZonesSuite{})

func (s *ZonesSuite) SetUpSuite(c *C) {
}

func (s *LabelsSuite) TestZonesLoad(c *C) {
	zs := new(Zones)
	err := zs.LoadZonesConfig("testdata/zones.json")
	c.Assert(err, IsNil)

	z, ok := zs.zones["z.example.com"]
	c.Assert(ok, Equals, true)
	c.Check(z.Name, Equals, "z.example.com")
	c.Check(z.Options.Ttl, Equals, 300) // default TTL

	z, ok = zs.zones["x.example.com"]
	c.Log(spew.Sdump(z))
	c.Assert(ok, Equals, true)
	c.Check(z.Name, Equals, "x.example.com")
	c.Check(z.Options.Ttl, Equals, 120) // Configured TTL

	/*

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

	*/

}
