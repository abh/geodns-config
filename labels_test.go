package main

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
	s.Labels.SetNode("label1", labelNode{Name: "foo", Ip: net.ParseIP("10.0.0.1")})
	c.Assert(s.Labels.Count(), Equals, 1)
	label := s.Labels.Get("label1")
	c.Assert(label, NotNil)
	c.Assert(label.Name, Equals, "label1")
}

func (s *LabelsSuite) TestLoad(c *C) {
	s.Labels.Clear()
	err := s.Labels.LoadFile("testdata/labels.json")
	c.Assert(err, IsNil)
}
