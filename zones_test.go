package dnsconfig

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
	c.Check(z.Options.Targeting, Equals, "@ country")

	c.Check(z.Ns, DeepEquals, []string{"ns1.example.com", "ns2.example.com"})

}
