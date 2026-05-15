package dnsconfig

import (
	"net/netip"

	. "launchpad.net/gocheck"
)

type NodesSuite struct {
	Nodes Nodes
}

var _ = Suite(&NodesSuite{})

func (s *NodesSuite) SetUpSuite(c *C) {
	s.Nodes = NewNodes()
}

func (s *NodesSuite) TestNodes(c *C) {
	s.Nodes.Clear()
	s.Nodes.Set("foo", Node{IP: netip.MustParseAddr("10.0.0.1"), Active: true})
	c.Assert(s.Nodes.Count(), Equals, 1)

	node := s.Nodes.Get("foo")
	c.Assert(node, NotNil)
	c.Assert(node.Name, Equals, "foo")
}

func (s *NodesSuite) TestLoad(c *C) {
	s.Nodes.Clear()
	err := s.Nodes.LoadFile("testdata/nodes-small.json")
	c.Assert(err, IsNil)

	node := s.Nodes.Get("edge01.lax")
	c.Assert(node, NotNil)
	c.Assert(node.IP.String(), Equals, "108.161.187.3")
	c.Assert(node.Active, Equals, true)

	node = s.Nodes.Get("edge01.sea")
	c.Assert(node, NotNil)
	c.Assert(node.Active, Equals, false)

	node = s.Nodes.Get("cname-chain")
	c.Assert(node, NotNil)
	c.Assert(node.Active, Equals, true)
	c.Assert(node.Cname, Equals, "hello.example.com")

	// not sure this is the appropriate data, but it's what's implemented so make it explicit
	c.Check(node.IP.IsValid(), Equals, false)
}
