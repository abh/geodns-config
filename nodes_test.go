package dnsconfig

import (
	"net/netip"
	"testing"
)

func TestNodesBasic(t *testing.T) {
	nodes := NewNodes()
	nodes.Set("foo", Node{IP: netip.MustParseAddr("10.0.0.1"), Active: true})
	if got := nodes.Count(); got != 1 {
		t.Fatalf("Count = %d, want 1", got)
	}

	node := nodes.Get("foo")
	if node == nil {
		t.Fatal("Get(foo) returned nil")
	}
	if node.Name != "foo" {
		t.Errorf("Name = %q, want %q", node.Name, "foo")
	}
}

func TestNodesLoad(t *testing.T) {
	nodes := NewNodes()
	if err := nodes.LoadFile("testdata/nodes-small.json"); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	node := nodes.Get("edge01.lax")
	if node == nil {
		t.Fatal("edge01.lax: nil")
	}
	if got := node.IP.String(); got != "108.161.187.3" {
		t.Errorf("edge01.lax IP = %q, want %q", got, "108.161.187.3")
	}
	if !node.Active {
		t.Error("edge01.lax Active = false, want true")
	}

	node = nodes.Get("edge01.sea")
	if node == nil {
		t.Fatal("edge01.sea: nil")
	}
	if node.Active {
		t.Error("edge01.sea Active = true, want false")
	}

	node = nodes.Get("cname-chain")
	if node == nil {
		t.Fatal("cname-chain: nil")
	}
	if !node.Active {
		t.Error("cname-chain Active = false, want true")
	}
	if node.Cname != "hello.example.com" {
		t.Errorf("cname-chain Cname = %q, want %q", node.Cname, "hello.example.com")
	}
	if node.IP.IsValid() {
		t.Errorf("cname-chain IP unexpectedly valid: %v", node.IP)
	}
}
