package DnsConfig 1.0;
use 5.12.0;
use Moose;
use Data::Dump qw(pp);
use Clone ();
use JSON  ();
use File::Slurp qw(read_file write_file);

my $json = JSON->new->canonical;

sub log { shift->config->log(@_) }

has 'dns' => (
    isa     => 'HashRef',
    is      => 'rw',
    default => sub { _base_data() }
);

sub _base_data() {
    return {
        data      => {"" => {ns => {"b2.bitnames.com" => "", "b3.bitnames.com" => ""}}},
        max_hosts => 2,
        ttl       => 60,
    };
}

has 'config' => (
    isa      => 'GeoConfig',
    is       => 'ro',
    required => 1,
    weak_ref => 1,
);


sub _setup_geo_rules {
    my $self = shift;
    my ($dns_name, $data) = (shift, shift);
    my $dns = {};

    my @pops = ref $data eq 'ARRAY' ? @$data : keys %$data;

    #Test::More::diag("geo: ", pp($dns_name), " POPS: ", pp(\@pops), " DATA: ", pp($data));

    my %ips;
    if (ref $data eq 'HASH') {
        for my $pop (keys %$data) {
            if ($data->{$pop}) {
                $ips{$pop} = $data->{$pop};
            }
        }
    }

    for my $pop (@pops) {
        my $geos = $self->config->pop_geo($pop);

        if (!$geos) {
            $self->log->warn("$pop not configured in geo.json");
            next;
        }
        for my $geo (@$geos) {
            my $pop_ip = $self->config->nodes->node_ip($pop) || $ips{$pop};
            unless ($pop_ip) {
                $self->log->warn("Unknown pop [$pop]");
                next;
            }
            if ($self->config->monitor && $self->config->monitor->outage($pop_ip)) {
                $self->log->info("$pop ($pop_ip) has a current outage");
                next;
            }

            my $geo_name = $dns_name;
            if ($geo ne '@') {
                $geo_name .= ".$geo";
            }
            $dns->{$geo_name} ||= {a => []};
            my $ip =  $ips{$pop} || $self->config->nodes->node_ip($pop);
            push @{$dns->{$geo_name}->{a}}, [$ip];
        }
    }
    unless ($dns->{$dns_name} && $dns->{$dns_name}->{a} && @{$dns->{$dns_name}->{a}}) {
        warn "$dns_name does not have any default records";
    }

    #say "DNS: ", pp($dns);

    return $dns;
}

sub add_geo_rules {
    my ($self, $name, $data) = @_;
    my $add = $self->_setup_geo_rules($name, $data);
    $self->dns->{data} = hash_merge($self->dns->{data}, $add);
}

sub setup_groups {
    my $self = shift;
    while (my ($group_name, $group_data) = each %{$self->config->groups}) {
        $self->log->debug("Configuring group $group_name");
        my $dns_name = "_" . $group_name;
        $self->add_geo_rules($dns_name, $group_data);
    }
    return 1;
}

sub setup_labels {
    my $self   = shift;
    my $groups = $self->config->groups;
    while (my ($domain_name, $domain_data) = each %{$self->config->labels->all}) {
        pp($domain_name, $domain_data);
        if ($self->dns->{data}->{$domain_name}) {
            warn "$domain_name configured twice (groups and domains?)";
            next;
        }
        if (my $alias = $domain_data->{group}) {
            if (!$groups->{$alias}) {
                warn "No group $alias configured";
                next;
            }
            # warn "adding alias for $domain_name / $alias";
            $self->dns->{data}->{$domain_name}->{alias} = "_$alias";
        }
        else {
            $self->add_geo_rules($domain_name, $domain_data);
        }
    }
    1;
}

sub setup_data {
    my $self = shift;
    #say "setting up data:", pp($self->config->labels);
    $self->dns(_base_data());
    $self->setup_groups;
    $self->setup_labels;
}

sub write_dns {
    my ($self, $file) = @_;
    my $dns_config = $self->dns;
    #pp($dns_config);
    $self->config->dirty(0);
    write_file($file, $json->encode($dns_config));
}

# $h = hash_merge($h1, $h2);
# $h = hash_merge($h1, $h2, $h3, ...);
sub hash_merge {

    # Do a deep copy of arguments, to avoid sharing
    my @h = @{Clone::clone(\@_)};

    # Merge them from left to right
    my $h = shift @h;
    _hash_merge($h, $_) for @h;
    return $h;
}

sub _hash_merge {
    my ($h1, $h2) = @_;

    if ($h2->{'__FINAL__'}) {
        %$h1 = %$h2;
        delete $h1->{'__FINAL__'};
        return;
    }

    keys %$h2;    # reset iter
    while (my ($k, $v) = each %$h2) {
        if (ref($v) eq 'HASH' and ref($h1->{$k}) eq 'HASH') {
            _hash_merge($h1->{$k}, $v);
        }
        elsif (defined($v) and !ref($v) and (my $tmp = $v) eq '__KILL__') {
            delete $h1->{$k};
        }
        else {
            $h1->{$k} = $v;
        }
    }
}


1;
