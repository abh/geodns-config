package GeoDNS::Nodes::NetDNA;
use Moose;
extends 'GeoDNS::Nodes';
with 'GeoDNS::Data::NetDNA';

sub update {
    my $self = shift;
    $self->last_check(time);
    my $res = $self->api->get('nodes.json');
    if ($res->is_status_class(200)) {
        my $nodes = $res->json('/data/nodes');
        my %nodes = map { %$_ } @$nodes;
        $self->{data} = \%nodes;
        $self->ready(1);
        $self->dirty(1);
    }
}

1;
