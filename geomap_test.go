package dnsconfig

import "testing"

func TestGeoMapLoad(t *testing.T) {
	g := NewGeoMap()
	if err := g.LoadFile("testdata/geomap.json"); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	if got := g.geomap["*.ams*"][0].target; got != "fr" {
		t.Errorf(`geomap["*.ams*"][0].target = %q, want %q`, got, "fr")
	}
	if got := g.geomap["*.ams*"][1].target; got != "nl" {
		t.Errorf(`geomap["*.ams*"][1].target = %q, want %q`, got, "nl")
	}
	if got := g.geomap["*.ams*"][2].target; got != "europe" {
		t.Errorf(`geomap["*.ams*"][2].target = %q, want %q`, got, "europe")
	}

	if got := g.GetNodeGeos("test.sea")[0].target; got != "@" {
		t.Errorf(`GetNodeGeos("test.sea")[0].target = %q, want "@"`, got)
	}

	if got := g.geomap["*.lhr"][1].weight; got != 1000 {
		t.Errorf(`geomap["*.lhr"][1].weight = %d, want 1000`, got)
	}

	if got := g.GetNodeGeos("flex04.ams04")[0].weight; got != 1 {
		t.Errorf(`GetNodeGeos("flex04.ams04")[0].weight = %d, want 1`, got)
	}
	if got := g.GetNodeGeos("flex04.ams04")[0].target; got != "europe" {
		t.Errorf(`GetNodeGeos("flex04.ams04")[0].target = %q, want "europe"`, got)
	}

	if got := g.GetNodeGeos("x123.lhr")[1].weight; got != 1000 {
		t.Errorf(`GetNodeGeos("x123.lhr")[1].weight = %d, want 1000`, got)
	}
	if got := g.GetNodeGeos("x123.lhr")[1].target; got != "europe" {
		t.Errorf(`GetNodeGeos("x123.lhr")[1].target = %q, want "europe"`, got)
	}
	if got := g.GetNodeGeos("x123.lhr")[0].weight; got != 100 {
		t.Errorf(`GetNodeGeos("x123.lhr")[0].weight = %d, want 100`, got)
	}
}
