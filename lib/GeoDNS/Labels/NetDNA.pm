package GeoDNS::Labels::NetDNA;
use Moose;
extends 'GeoDNS::Labels';
with 'GeoDNS::Data::NetDNA';

sub update {
    my $self = shift;
    $self->last_check(time);
    $self->api->get(
        'flexnodes.json',
        sub {
            my $res = shift;
            if ($res->is_status_class(200)) {
                my $labels = $res->json('/data/nodes');
                $self->{data} = $labels;
                $self->ready(1);
                $self->dirty(1);
            }
        }
    );
}

sub all {
    my $self = shift;
    return \%{$self->{data}};
}

1;
