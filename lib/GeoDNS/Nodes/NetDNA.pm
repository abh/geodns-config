package GeoDNS::Nodes::NetDNA;
use Moose;
extends 'GeoDNS::Nodes';
with 'GeoDNS::Data::NetDNA', 'GeoDNS::Log::Role';

sub update {
    my $self = shift;
    $self->last_check(time);
    my $delay = Mojo::IOLoop->delay;
    $delay->begin;
    $self->api->get(
        'nodes.json',
        sub {
            my ($res) = @_;
            if ($res->is_status_class(200)) {
                my $nodes = $res->json('/data/nodes');
                my %nodes = map {%$_} @$nodes;
                for my $name (keys %nodes) {
                    $self->log->debug("data for $name: ", Data::Dump::pp($nodes{$name}));
                    unless ($nodes{$name}->{ip}) {
                        $self->log->info("node $name doesn't have an IP address, skipping");
                        delete $nodes{$name};
                        next;
                    }
                    $nodes{$name}->{name} = $name;
                    $nodes{$name} = GeoDNS::Node->new($nodes{$name});
                }
                $self->nodes(\%nodes);
                $self->ready(1);
                $self->dirty(1);
            }
            else {
                $self->log->warn('nodes.json failed: ' . $res->code);
                $self->log->info('nodes.json return',  Data::Dump::pp($res));
            }
            $delay->end;

        }
    );
}

1;
