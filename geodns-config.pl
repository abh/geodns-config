#!/usr/bin/env perl
use v5.12.0;
use strict;
use JSON qw(decode_json encode_json);
use File::Slurp qw(read_file write_file);
use Data::Dump qw(pp);

# TODO
#  - failsafe handling all servers in a group being "out".
#    - maybe mark some pops as "can't go down?" (anycast)


my $pops      = decode_json(read_file("config/pops.json"));
my $groups    = decode_json(read_file("config/groups.json"));
my $domains   = decode_json(read_file("config/domains.json"));
my $geoconfig = decode_json(read_file("config/geo.json"));
my $outages   = decode_json(read_file("data/outages.json"));

pp($pops);

my $dns_config = {
    data      => {"" => {ns => {"b2.bitnames.com" => "", "b3.bitnames.com" => ""}}},
    max_hosts => 2,
    ttl       => 60,
};

sub _geo_rules {
    my ($dns_name, $data) = (shift, shift);
    my $dns = $dns_config->{data};

    #say "geo: ", pp($geoconfig), " POPS: ", pp($pops), " DATA: ", pp($data);

    my @pops = ref $data eq 'ARRAY' ? @$data : keys %$data;
    my %ips = %$pops;
    if (ref $data eq 'HASH') {
        for my $pop (keys %$data) {
            if ($data->{$pop}) {
                $ips{$pop} = $data->{$pop};
            }
        }
    }
    say "IPS: ", pp(\%ips);

    for my $pop (@pops) {
        my $geos = $geoconfig->{$pop};
        if (!$geos) {
            warn "$pop not configured in geo.json";
            next;
        }
        for my $geo (@$geos) {
            my $pop_ip = $pops->{$pop} || $ips{$pop};
            unless ($pop_ip) {
                warn "Unknown pop [$pop]";
                next;
            }
            if ($outages->{$pop_ip}) {
                say "$pop ($pop_ip) has a current outage";
                next;
            }

            my $geo_name = $dns_name;
            if ($geo ne '@') {
                $geo_name .= ".$geo";
            }
            $dns->{$geo_name} ||= {a => []};
            push @{$dns->{$geo_name}->{a}}, [$ips{$pop}];
        }
    }
    unless ($dns->{$dns_name} && $dns->{$dns_name}->{a} && @{$dns->{$dns_name}->{a}}) {
        warn "$dns_name does not have any default records";
    }
}

while (my ($group_name, $group_data) = each %$groups) {
    say $group_name;
    my $dns_name = "_" . $group_name;
    _geo_rules($dns_name, $group_data);
}

my $dns = $dns_config->{data};

while (my ($domain_name, $domain_data) = each %$domains) {
    pp($domain_name, $domain_data);
    if ($dns_config->{data}->{$domain_name}) {
        die "$domain_name configured twice (groups and domains?)";
    }
    if (my $alias = $domain_data->{alias}) {
        if (!$groups->{$alias}) {
            die "No group $alias configured";
            next;
        }
        $dns->{$domain_name}->{alias} = "_$alias";
    }
    else {
        _geo_rules($domain_name, $domain_data);
    }
}

pp($dns_config);
write_file("dns/g.develooper.org.json", encode_json($dns_config));
