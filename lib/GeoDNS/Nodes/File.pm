package GeoDNS::Nodes::File;
use Moose;
extends 'GeoDNS::Nodes';
with 'GeoDNS::JsonFile', 'GeoDNS::Log::Role';

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

1;
