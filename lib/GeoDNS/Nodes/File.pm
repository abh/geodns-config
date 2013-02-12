package GeoDNS::Nodes::File;
use Moose;
extends 'GeoDNS::Nodes';
with 'GeoDNS::JsonFile', 'GeoConfig::Log';

has 'name' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'nodes',
);

sub update {
    my $self = shift;
    if ($self->_refresh_data('data', $self->file)) {
        for my $node (keys %{$self->{data}}) {
            unless (ref $self->{data}->{$node}) {
                $self->{data}->{$node} = {ip => $self->{data}->{$node}, active => 1};
            }
        }
        return $self->ready(1);
    }
    return 0;
}

sub node_ip {
    my ($self, $node) = @_;
    $self->update();

    #Test::More::diag("Nodes, fetching '$node': " .$self->{data}->{$node} ." / ". Data::Dump::pp($self->{data}));
    return $self->{data}->{$node}->{ip};
}

sub set_ip {
    my ($self, $node, $ip, $active) = @_;
    unless (defined $ip) {
        $active = 1;
    }
    return $self->{data}->{$node} = {ip => $ip, active => $active};
}

sub all {
    my $self = shift;
    return \%{$self->{data}};
}

1;
