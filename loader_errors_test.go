package dnsconfig

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoaderErrors verifies that malformed JSON inputs return errors rather
// than panicking. Each case exercises a type-assertion site that used to
// panic before the typed-struct refactor.
func TestLoaderErrors(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		load     func(path string) error
	}{
		{
			"nodes: value is not an object",
			`{"foo": 42}`,
			func(p string) error { ns := NewNodes(); return ns.LoadFile(p) },
		},
		{
			"nodes: missing ip and cname",
			`{"foo": {"active": 1}}`,
			func(p string) error { ns := NewNodes(); return ns.LoadFile(p) },
		},
		{
			"nodes: invalid ip string",
			`{"foo": {"ip": "not-an-ip", "active": 1}}`,
			func(p string) error { ns := NewNodes(); return ns.LoadFile(p) },
		},
		{
			"geomap: target element is not a string",
			`{"foo": [42]}`,
			func(p string) error { gm := NewGeoMap(); return gm.LoadFile(p) },
		},
		{
			"geomap: bad weight",
			`{"foo": ["us=abc"]}`,
			func(p string) error { gm := NewGeoMap(); return gm.LoadFile(p) },
		},
		{
			"zones: contact is not a string",
			`{"z.example": {"contact": 42}}`,
			func(p string) error { zs := new(Zones); return zs.LoadZonesConfig(p) },
		},
		{
			"labels: entry value is a number",
			`{"l": {"node": 42}}`,
			func(p string) error { ls := NewLabels(); return ls.LoadFile(p) },
		},
		{
			"labels: invalid ip string",
			`{"l": {"node": "not-an-ip"}}`,
			func(p string) error { ls := NewLabels(); return ls.LoadFile(p) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "test.json")
			if err := os.WriteFile(p, []byte(tc.contents), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
			if err := tc.load(p); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
