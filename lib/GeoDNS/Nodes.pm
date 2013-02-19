package GeoDNS::Nodes;
use Moose;
use GeoDNS::Node;

has nodes => (
    isa => 'HashRef[GeoDNS::Node]',
    is  => 'rw',
    default => sub { {} }
);

sub node_ip {
    my ($self, $node) = @_;
    return unless $self->nodes->{$node};
    return $self->nodes->{$node}->{ip};
}

sub set_ip {
    my ($self, $node, $ip, $active) = @_;
    unless (defined $ip) {
        $active = 1;
    }
    return $self->nodes->{$node} = {ip => $ip, active => $active};
}

sub all {
    my $self = shift;
    return \%{$self->nodes};
}


1;
