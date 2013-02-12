package GeoDNS::Nodes;
use Moose;

sub node_ip {
    my ($self, $node) = @_;
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
