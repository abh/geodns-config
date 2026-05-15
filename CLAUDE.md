# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`geodns-config` is a Go library and CLI tool (`dnsconfig`) that converts compact JSON configuration into [GeoDNS](http://geo.bitnames.com/) zone files. It reads zones, nodes, geomap, and labels configs and outputs GeoDNS-format JSON zone files.

## Commands

```bash
# Run all tests
go test ./...

# Run a single test
go test -run TestName ./...

# Build the dnsconfig binary (from dnsconfig/ directory)
cd dnsconfig && ./build/darwin    # macOS
cd dnsconfig && ./build/linux     # Linux

# Run the tool (from repo root)
./dnsconfig/dnsconfig -config config/zones.json -output dns/
```

## Architecture

**Package layout:**
- Root package `dnsconfig` — the reusable library
- `dnsconfig/` directory — the `main` package CLI binary

**Core data flow:** `zones.json` → `Zone.LoadConfig()` → `Zone.BuildJSON()` → GeoDNS JSON files in `dns/`

**Primary types** (each guards its map with a `sync.Mutex`):
- `Zones` (`zones.go`) — `LoadZonesConfig()` reads `zones.json` and instantiates a `Zone` per entry; each `Zone.LoadConfig()` then calls the three loaders below
- `Nodes` (`nodes.go`) — `LoadFile()` reads `nodes.json`, maps node names to IP/CNAME/active state (IPs are `netip.Addr`)
- `Labels` (`labels.go`) — `LoadFile()` reads `labels.json`, maps hostnames to node references; supports group aliases and per-node IP overrides
- `GeoMap` (`geomap.go`) — `LoadFile()` reads `geomap.json`, maps node name patterns to geo targets with weights

**Key helpers:**
- `jsonloader.go`: `jsonLoader()` wraps JSON decoding with line/column error reporting via `github.com/abh/errorutil`
- `utils.go`: `matchWildcard()` converts wildcard patterns (e.g. `*.ams`) to regex; also handles raw regex when the key starts with `^` or ends with `$`
- `dns.go`: `BuildZone()` joins Nodes + Labels + GeoMap into a `zoneJson` struct; `BuildJSON()` marshals it

**Wildcard/regex matching** is used in both geomap keys (node name patterns) and label node keys — both delegate to `matchWildcard()`.

**Config file paths** in `zones.json` are relative to the zones file location; `absPath()` in `zones.go` resolves them.

## Test data

`testdata/` contains sample JSON files used by tests. `t/config-test/` holds an integration test config. `config/` and `dns/` are live NTP Pool configuration and output.
