package GeoDNS::Labels::File;
use Moose;
extends 'GeoDNS::Labels';
with 'GeoDNS::JsonFile', 'GeoConfig::Log';

has 'name' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'labels',
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
    if ($self->_refresh_data('labels', $self->file)) {
        return 1;
    }
    return 0;
}

sub all {
    my $self = shift;
    return \%{$self->{labels}};
}

1;
