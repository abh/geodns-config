# GeoDNS Config Tool

The `dnsconfig` tool helps create configuration/zone files for the
[GeoDNS](http://geo.bitnames.com/) server.

## Command line options

* -config: name of the zones config file. Defaults to `./config/zones.json`.
* -output: name of output directory. Defaults to `./dns/`.

## Configuration files

### Zones

By default `config/zones.json`. You can specify another with the `-config` parameter.

    {
        "some.example.com": {
                "contact": "noc@example.com",
                "ttl":     120,
                "max_hosts": 2,
                "ns":      ["ns1.example.com", "ns2.example.com"],
                "labels":  "labels.json",
                "nodes":   "nodes.json",
                "geomap":  "geomap.json"
        }
    }

The options are:

* contact: Used for the SOA contact field.
* ttl: Time-to-live configuration for the DNS replies (in seconds).
* max_hosts: Maximum number of IPs to return in each reply.
* ns: List of nameservers for the zone.
* labels, nodes and geomap: Filename for data configuration (see below). The filenames are relative to the location of the zone configuration.

Multiple zones can be specified in the file.

### Nodes

List of named "nodes" (servers). The name is only used as an identifier to
match data in the other configuration files, not in replies. Setting `active`
to zero will remove this node from the output data.

The IP address is a default to be used when the server is used in a DNS reply,
it can be overridden in the labels file.

    {
     "edge01.any": { "ip": "10.0.0.1", "active": 1 },
     "edge04.any": { "ip": "10.0.0.4", "active": 1 },

     "edge01.lax": { "ip": "10.0.1.1", "active": 1 },
     "edge01.sea": { "ip": "10.0.2.1", "active": 0 }
    }

### Geomap

A geomap maps the nodes to "targets" (countries and continents). Each node has
a list of targets it will "match".

The special "@" target is the default target if there are no more specific
matches.

An equal sign followed by a number specifies a "weight", higher weights will be
returned in answers more often. The default weight is 100.

The "key" in the data structure can have wildcards ("*") matching any non-dot
character. To match "foo.bar" you can use "*.bar", "foo.*" or "*.*".

    {
        "*.any": [ "@" ],
        "*.sin": [ "sg", "th", "id", "my" ],
        "*.ams": [ "europe", "nl", "fr" ],
        "*.lhr": [ "europe=1000", "uk" ],
        "*.sea": [ "us" ],
        "flex04.ams04": [ "europe=1" ]
    }

### Labels

Labels are 'host names' in the zone. The value for each key is a hash with node
names (must match an entry in the nodes config) and an optional IP override.

The override can also be another hash with the elements 'active' (defaults to true)
and 'ip' (optional). 'active' can be specified as true, 1, false or 0.

    {
        "some.example":  {
            "edge01.any": "",
            "flex01.sin": ""
        },
        "alias.example": {
            "group": "some.example"
        },
        "another.test": {
            "edge01.any": "10.1.1.10",
            "flex01.sin": "10.20.1.10",
            "edge01.lhr": ""
        },
        "zone4": {
            "edge01.sea": { "active": true, ip: "10.1.2.3" },
            "edge01.any": { "active": 0 }
        }
    }

## Copyright

Copyright 2013 Ask Bjørn Hansen