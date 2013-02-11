package GeoDNS::Nodes::File;
use Moose;
extends 'GeoDNS::Nodes';
with 'GeoDNS::JsonFile', 'GeoConfig::Log';

has 'name' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'nodes',
);

has 'file' => (
    isa     => 'Str',
    is      => 'rw',
    lazy    => 1,
    default => sub { shift->name . '.json' },
);

sub check {

    # run update if appropriate
    my $self = shift;
    return $self->update;
}

sub update {
    my $self = shift;
    if ($self->_refresh_data('nodes', $self->file)) {
        for my $node (keys %{$self->{nodes}}) {
            unless (ref $self->{nodes}->{$node}) {
                $self->{nodes}->{$node} = {ip => $self->{nodes}->{$node}, active => 1};
            }
        }
        return 1;
    }
    return 0;
}

sub node_ip {
    my ($self, $node) = @_;
    $self->update();

#Test::More::diag("Nodes, fetching '$node': " .$self->{nodes}->{$node} ." / ". Data::Dump::pp($self->{nodes}));
    return $self->{nodes}->{$node}->{ip};
}

sub set_ip {
    my ($self, $node, $ip, $active) = @_;
    unless (defined $ip) {
        $active = 1;
    }
    return $self->{nodes}->{$node} = {ip => $ip, active => $active};
}

sub all {
    my $self = shift;
    return \%{$self->{nodes}};
}

1;
