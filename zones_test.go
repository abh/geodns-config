package dnsconfig

import (
	"reflect"
	"testing"
)

func TestZonesLoad(t *testing.T) {
	zs := new(Zones)
	if err := zs.LoadZonesConfig("testdata/zones.json"); err != nil {
		t.Fatalf("LoadZonesConfig: %v", err)
	}

	z, ok := zs.zones["z.example.com"]
	if !ok {
		t.Fatal("z.example.com missing")
	}
	if z.Name != "z.example.com" {
		t.Errorf("Name = %q, want %q", z.Name, "z.example.com")
	}
	if z.Options.Ttl != 300 {
		t.Errorf("Ttl = %d, want default 300", z.Options.Ttl)
	}

	z, ok = zs.zones["x.example.com"]
	if !ok {
		t.Fatal("x.example.com missing")
	}
	if z.Name != "x.example.com" {
		t.Errorf("Name = %q, want %q", z.Name, "x.example.com")
	}
	if z.Options.Ttl != 120 {
		t.Errorf("Ttl = %d, want 120", z.Options.Ttl)
	}
	if z.Options.Targeting != "@ country" {
		t.Errorf("Targeting = %q, want %q", z.Options.Targeting, "@ country")
	}

	wantNs := []string{"ns1.example.com", "ns2.example.com"}
	if !reflect.DeepEqual(z.Ns, wantNs) {
		t.Errorf("Ns = %v, want %v", z.Ns, wantNs)
	}
}
