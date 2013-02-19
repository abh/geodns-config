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
    if ($self->_refresh_data($self->file,
        sub {
            my $data = shift;
            my %nodes;
            for my $name (keys %$data) {
                my $d = $data->{$name};
                unless (ref $d) {
                    $d = {ip => $data->{$name}, active => 1};
                }
                $d->{name} = $name;
                #warn "D: ", Data::Dump::pp($d);
                my $node = GeoDNS::Node->new($d);
                $nodes{$name} = $node;
            }
            $self->nodes(\%nodes);
            #warn "Got nodes: ", Data::Dump::pp($self->nodes), " from ", Data::Dump::pp(\%nodes);
            return $self->ready(1);
        }
      ))
    {
        return 1;
    }
    return 0;
}

1;
