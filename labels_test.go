package dnsconfig

import (
	"net/netip"
	"testing"
)

func TestLabelsBasic(t *testing.T) {
	labels := NewLabels()
	labels.SetNode("label1", labelNode{Name: "foo", IP: netip.MustParseAddr("10.0.0.1")})
	if got := labels.Count(); got != 1 {
		t.Fatalf("Count = %d, want 1", got)
	}
	label := labels.Get("label1")
	if label == nil {
		t.Fatal("Get(label1) returned nil")
	}
	if label.Name != "label1" {
		t.Errorf("Name = %q, want %q", label.Name, "label1")
	}
}

func TestLabelsLoad(t *testing.T) {
	labels := NewLabels()
	if err := labels.LoadFile("testdata/labels.json"); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	if got := labels.Get("zone1.example").GroupName; got != "edge1-global" {
		t.Errorf("zone1.example GroupName = %q, want %q", got, "edge1-global")
	}

	if got := labels.Get("zone2.example").GetNode("edge01.any").Name; got != "edge01.any" {
		t.Errorf("zone2.example edge01.any Name = %q", got)
	}
	if got := labels.Get("zone2.example").GetNode("edge01.any").IP.String(); got != "10.1.1.10" {
		t.Errorf("zone2.example edge01.any IP = %q, want %q", got, "10.1.1.10")
	}

	if ip := labels.Get("zone3.example").GetNode("edge01.any").IP; ip.IsValid() {
		t.Errorf("zone3.example edge01.any IP unexpectedly valid: %v", ip)
	}
	if !labels.Get("zone3.example").GetNode("edge01.any").Active {
		t.Error("zone3.example edge01.any Active = false, want true")
	}
	if got := labels.Get("zone3.example").GetNode("cname-one").Cname; got != "one-override.example.com" {
		t.Errorf("zone3.example cname-one Cname = %q", got)
	}
	if !labels.Get("zone3.example").GetNode("cname-one").Active {
		t.Error("zone3.example cname-one Active = false, want true")
	}

	if !labels.Get("zone4").GetNode("edge01.any").Active {
		t.Error("zone4 edge01.any Active = false, want true")
	}
	if labels.Get("zone4").GetNode("edge01.jfk").Active {
		t.Error("zone4 edge01.jfk Active = true, want false")
	}

	node := labels.Get("match").GetNode("edge01.jfk")
	if node == nil {
		t.Fatal("match edge01.jfk: nil")
	}
	if node.Name != "edge01.jfk" {
		t.Errorf("match Name = %q", node.Name)
	}
	if !node.Active {
		t.Error("match Active = false, want true")
	}

	node = labels.Get("match-inactive").GetNode("edge01.jfk")
	if node == nil {
		t.Fatal("match-inactive edge01.jfk: nil")
	}
	if node.Name != "edge01.jfk" {
		t.Errorf("match-inactive Name = %q", node.Name)
	}
	if node.Active {
		t.Error("match-inactive Active = true, want false")
	}
}
