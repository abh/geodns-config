package dnsconfig

import (
	. "launchpad.net/gocheck"
	"net"
)

type LabelsSuite struct {
	Labels Labels
}

var _ = Suite(&LabelsSuite{})

func (s *LabelsSuite) SetUpSuite(c *C) {
	s.Labels = NewLabels()
}

func (s *LabelsSuite) TestLabels(c *C) {
	s.Labels.Clear()
	s.Labels.SetNode("label1", labelNode{Name: "foo", IP: net.ParseIP("10.0.0.1")})
	c.Assert(s.Labels.Count(), Equals, 1)
	label := s.Labels.Get("label1")
	c.Assert(label, NotNil)
	c.Assert(label.Name, Equals, "label1")
}

func (s *LabelsSuite) TestLoad(c *C) {
	s.Labels.Clear()
	err := s.Labels.LoadFile("testdata/labels.json")
	c.Assert(err, IsNil)

	c.Assert(s.Labels.Get("zone1.example").GroupName, Equals, "edge1-global")
	c.Assert(s.Labels.Get("zone2.example").GetNode("edge01.any").Name, Equals, "edge01.any")
	c.Assert(s.Labels.Get("zone2.example").GetNode("edge01.any").IP.String(), Equals, "10.1.1.10")
	c.Assert(s.Labels.Get("zone3.example").GetNode("edge01.any").IP, IsNil)
	c.Assert(s.Labels.Get("zone3.example").GetNode("edge01.any").Active, Equals, true)

	c.Assert(s.Labels.Get("zone4").GetNode("edge01.any").Active, Equals, true)
	c.Assert(s.Labels.Get("zone4").GetNode("edge01.jfk").Active, Equals, false)

	// wildcards
	node := s.Labels.Get("match").GetNode("edge01.jfk")
	c.Assert(node, NotNil)
	c.Check(node.Name, Equals, "edge01.jfk")
	c.Check(node.Active, Equals, true)

	node = s.Labels.Get("match-inactive").GetNode("edge01.jfk")
	c.Assert(node, NotNil)
	c.Check(node.Name, Equals, "edge01.jfk")
	c.Check(node.Active, Equals, false)

}
