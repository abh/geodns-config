package GeoConfig 1.0;
use v5.12.0;
use Moose;
extends 'GeoConfig::Config';
use Data::Dump qw(pp);
use DnsConfig;

has 'dns' => (
    isa        => 'DnsConfig',
    is         => 'ro',
    lazy_build => 1,
);

sub _build_dns {
    my $self = shift;
    return DnsConfig->new(config => $self);
}

1;
