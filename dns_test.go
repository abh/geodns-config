package dnsconfig

import (
	"reflect"
	"testing"
)

func newTestZone(t *testing.T) *Zone {
	t.Helper()
	z := new(Zone)
	z.Name = "example.com"
	z.Options.Ttl = 25
	if err := z.Labels.LoadFile("testdata/labels.json"); err != nil {
		t.Fatalf("Labels.LoadFile: %v", err)
	}
	if err := z.GeoMap.LoadFile("testdata/geomap.json"); err != nil {
		t.Fatalf("GeoMap.LoadFile: %v", err)
	}
	if err := z.Nodes.LoadFile("testdata/nodes.json"); err != nil {
		t.Fatalf("Nodes.LoadFile: %v", err)
	}
	return z
}

func TestDnsLoad(t *testing.T) {
	z := newTestZone(t)
	zd, err := z.BuildZone()
	if err != nil {
		t.Fatalf("BuildZone: %v", err)
	}

	if zd.Ttl != 25 {
		t.Errorf("Ttl = %d, want 25", zd.Ttl)
	}

	t1, ok := zd.Data["zone2.example"]
	if !ok {
		t.Fatal("zone2.example missing")
	}
	// IP override, default weight
	want := []interface{}{"10.1.1.10", 100}
	if got := t1.A[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("zone2.example A[0] = %v, want %v", got, want)
	}

	t1, ok = zd.Data["zone2.example.europe"]
	if !ok {
		t.Fatal("zone2.example.europe missing")
	}
	// use default IP from nodes.json and weight override
	want = []interface{}{"10.0.5.1", 1000}
	if got := t1.A[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("zone2.example.europe A[0] = %v, want %v", got, want)
	}

	if _, ok := zd.Data["zone4"]; !ok {
		t.Error("zone4 missing")
	}

	// edge01.sea is inactive and edge01.jfk disabled for this label
	if _, ok := zd.Data["zone4.us"]; ok {
		t.Error("zone4.us should be missing (both us nodes inactive/disabled)")
	}

	t1, ok = zd.Data["any-only"]
	if !ok {
		t.Fatal("any-only missing")
	}
	want = []interface{}{"10.0.0.1", 100}
	if got := t1.A[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("any-only A[0] = %v, want %v", got, want)
	}
	want = []interface{}{"10.0.0.4", 100}
	if got := t1.A[1].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("any-only A[1] = %v, want %v", got, want)
	}

	t1, ok = zd.Data["any-only.north-america"]
	if !ok {
		t.Fatal("any-only.north-america missing")
	}
	want = []interface{}{"10.0.0.1", 10}
	if got := t1.A[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("any-only.north-america A[0] = %v, want %v", got, want)
	}
	want = []interface{}{"10.0.0.4", 10}
	if got := t1.A[1].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("any-only.north-america A[1] = %v, want %v", got, want)
	}

	t1, ok = zd.Data["any-alias"]
	if !ok {
		t.Fatal("any-alias missing")
	}
	if t1.Alias != "any-only" {
		t.Errorf("any-alias Alias = %q, want %q", t1.Alias, "any-only")
	}

	js, err := zd.JSON()
	if err != nil {
		t.Errorf("JSON: %v", err)
	}
	if len(js) == 0 {
		t.Error("JSON returned empty string")
	}
}

func TestAaaa(t *testing.T) {
	z := newTestZone(t)
	zd, err := z.BuildZone()
	if err != nil {
		t.Fatalf("BuildZone: %v", err)
	}

	// IPv6-only node lands in Aaaa, not A
	t1, ok := zd.Data["v6-only"]
	if !ok {
		t.Fatal("v6-only missing")
	}
	if got := len(t1.A); got != 0 {
		t.Errorf("v6-only len(A) = %d, want 0", got)
	}
	want := []interface{}{"2001:db8::6", 100}
	if got := t1.Aaaa[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("v6-only Aaaa[0] = %v, want %v", got, want)
	}

	// Dual-stack label has both A and AAAA records
	t1, ok = zd.Data["dual-stack"]
	if !ok {
		t.Fatal("dual-stack missing")
	}
	want = []interface{}{"10.0.6.6", 100}
	if got := t1.A[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("dual-stack A[0] = %v, want %v", got, want)
	}
	want = []interface{}{"2001:db8::6", 100}
	if got := t1.Aaaa[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("dual-stack Aaaa[0] = %v, want %v", got, want)
	}

	// IP override with an IPv6 address routes to Aaaa even when the
	// node default is IPv4
	t1, ok = zd.Data["v6-override"]
	if !ok {
		t.Fatal("v6-override missing")
	}
	if got := len(t1.A); got != 0 {
		t.Errorf("v6-override len(A) = %d, want 0", got)
	}
	want = []interface{}{"2001:db8::101", 100}
	if got := t1.Aaaa[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("v6-override Aaaa[0] = %v, want %v", got, want)
	}
}

func TestCname(t *testing.T) {
	z := newTestZone(t)
	zd, err := z.BuildZone()
	if err != nil {
		t.Fatalf("BuildZone: %v", err)
	}

	t1, ok := zd.Data["zone3.example.dk"]
	if !ok {
		t.Fatal("zone3.example.dk missing")
	}
	want := []interface{}{"one-override.example.com", 100}
	if got := t1.Cname[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("zone3.example.dk Cname[0] = %v, want %v", got, want)
	}

	t1, ok = zd.Data["zone3.example.se"]
	if !ok {
		t.Fatal("zone3.example.se missing")
	}
	want = []interface{}{"two.example.com", 100}
	if got := t1.Cname[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("zone3.example.se Cname[0] = %v, want %v", got, want)
	}

	t1, ok = zd.Data["zone3.example.no"]
	if !ok {
		t.Fatal("zone3.example.no missing")
	}
	want = []interface{}{"one-override.example.com", 2}
	if got := t1.Cname[0].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("zone3.example.no Cname[0] = %v, want %v", got, want)
	}
	want = []interface{}{"two.example.com", 1}
	if got := t1.Cname[1].([]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("zone3.example.no Cname[1] = %v, want %v", got, want)
	}
}

func TestDnsSort(t *testing.T) {
	zd := zoneData{}
	l := new(zoneLabel)
	zd["test"] = l

	// Inputs chosen so numeric and lexicographic orders disagree:
	// lex order: 10.x, 192.x, 21.x, 9.x; numeric: 9, 10, 21, 192.
	l.A = jsonAddresses{
		[]interface{}{"192.0.2.1", 100},
		[]interface{}{"21.0.0.1", 200},
		[]interface{}{"9.0.0.1", 300},
		[]interface{}{"10.0.0.1", 400},
	}
	l.Aaaa = jsonAddresses{
		[]interface{}{"2001:db8::20", 100},
		[]interface{}{"2001:db8::3", 200},
	}
	l.Cname = jsonAddresses{
		[]interface{}{"b.example.com", 100},
		[]interface{}{"a.example.com", 200},
	}

	zd.sortRecords()

	wantA := jsonAddresses{
		[]interface{}{"9.0.0.1", 300},
		[]interface{}{"10.0.0.1", 400},
		[]interface{}{"21.0.0.1", 200},
		[]interface{}{"192.0.2.1", 100},
	}
	if !reflect.DeepEqual(l.A, wantA) {
		t.Errorf("A sorted = %v, want %v", l.A, wantA)
	}

	wantAaaa := jsonAddresses{
		[]interface{}{"2001:db8::3", 200},
		[]interface{}{"2001:db8::20", 100},
	}
	if !reflect.DeepEqual(l.Aaaa, wantAaaa) {
		t.Errorf("Aaaa sorted = %v, want %v", l.Aaaa, wantAaaa)
	}

	wantCname := jsonAddresses{
		[]interface{}{"a.example.com", 200},
		[]interface{}{"b.example.com", 100},
	}
	if !reflect.DeepEqual(l.Cname, wantCname) {
		t.Errorf("Cname sorted = %v, want %v", l.Cname, wantCname)
	}
}
