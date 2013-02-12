package GeoConfig 1.0;
use v5.12.0;
use Moose;
extends 'GeoConfig::Config';
use Data::Dump qw(pp);
use DnsConfig;

has 'domain_name' => (
    isa      => 'Str',
    is       => 'rw',
    required => 1,
);

has 'dns' => (
    isa        => 'DnsConfig',
    is         => 'ro',
    lazy_build => 1,
);

sub _build_dns {
    my $self = shift;
    return DnsConfig->new(config => $self);
}

has 'monitor' => (
    isa     => 'GeoDNS::Monitor',
    is      => 'ro',
    default => sub {
        require GeoDNS::Monitor::Manual;
        GeoDNS::Monitor::Manual->new;
    },
);

has 'nodes' => (
    isa     => 'GeoDNS::Nodes',
    is      => 'ro',
    lazy    => 1,
    default => sub {
        my $self = shift;
        require GeoDNS::Nodes::File;
        return GeoDNS::Nodes::File->new(
            file => join("/", $self->config_path, $self->domain_name . '.nodes.json'),
            name => $self->domain_name,
        );
    },
);

has 'labels' => (
    isa     => 'GeoDNS::Labels',
    is      => 'ro',
    lazy    => 1,
    default => sub {
        my $self = shift;
        require GeoDNS::Labels::File;
        return GeoDNS::Labels::File->new(
            file => join("/", $self->config_path, $self->domain_name . '.labels.json'),
            name => $self->domain_name,
        );
    },
);

sub ready {
    my $self = shift;
    return $self->labels->ready && $self->nodes->ready;
}

sub pop_geo {
    my $self = shift;
    my $pop = shift;
    my $geoconfig = $self->geoconfig;
    if (my $geo = $geoconfig->{$pop}) {
        return $geo;
    }
    my @wc = map {
        my $r   = $_;
        my $geo = $geoconfig->{$r};
        $r =~ s/\./\\./g;
        $r =~ s/\*/[^\.]+/g;
        my $re = qr/^$r$/;
        [$re => $geo]
    } grep { m/^\*/ or m/\.\*/ } keys %$geoconfig;

    for my $rule (@wc) {
        if ($pop =~ $rule->[0]) {
            return $rule->[1];
        }
    }

    return [];
}

1;
