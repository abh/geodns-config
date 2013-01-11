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
