package GeoConfig::Config;
use Moose;
with 'GeoDNS::JsonFile', 'GeoConfig::Log';
use Data::Dump qw(pp);

has 'config_path' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'config'
);

sub BUILD {
    my $self = shift;
    $self->refresh;
}

has 'groups' => (
    isa        => 'HashRef',
    is         => 'rw',
);

sub refresh_groups {
    my $self = shift;
    return $self->_refresh_data('groups', 'groups.json');
}

has 'geoconfig' => (
    isa        => 'HashRef',
    is         => 'rw',
);

sub refresh_geoconfig {
    my $self = shift;
    return $self->_refresh_data('geoconfig', 'geo.json');
}

after 'dirty' => sub {
    my $self = shift;

    if (defined $_[0] && $_[0] == 0) {
        $self->nodes->dirty(0);
        $self->labels->dirty(0);
    }
};

sub refresh {
    my $self = shift;
    $self->refresh_groups;
    $self->refresh_geoconfig;
    $self->nodes->check;
    $self->labels->check;
    if (!$self->dirty) {
        if (my $dirty = $self->nodes->dirty || $self->labels->dirty) {
            $self->dirty($dirty) if $dirty;
        }
    }
    return $self->dirty;
}

1;
