package GeoConfig::Config;
use Moose;
extends 'GeoConfig::Log';
use Data::Dump qw(pp);
use JSON qw(decode_json encode_json);
use File::Slurp qw(read_file write_file);

has 'last_read_timestamp' => (
    isa     => 'HashRef',
    is      => 'rw',
    default => sub { return {} }
);

has 'config_path' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'config'
);

has 'pops' => (
    isa        => 'HashRef',
    is         => 'rw',
    lazy_build => 1,
    builder    => 'refresh_pops'
);

sub refresh_pops {
    my $self = shift;
    return $self->_refresh_data('pops', 'pops.json');
}

has 'groups' => (
    isa        => 'HashRef',
    is         => 'rw',
    lazy_build => 1,
    builder    => 'refresh_groups'
);

sub refresh_groups {
    my $self = shift;
    return $self->_refresh_data('groups', 'groups.json');
}

has 'labels' => (
    isa        => 'HashRef',
    is         => 'rw',
    lazy_build => 1,
    builder    => 'refresh_labels'
);

sub refresh_labels {
    my $self = shift;
    return $self->_refresh_data('labels', 'labels.json');
}

has 'geoconfig' => (
    isa        => 'HashRef',
    is         => 'rw',
    lazy_build => 1,
    builder    => 'refresh_geoconfig'
);

sub refresh_geoconfig {
    my $self = shift;
    return $self->_refresh_data('geoconfig', 'geo.json');
}

sub refresh {
    my $self = shift;
    $self->refresh_pops;
    $self->refresh_groups;
    $self->refresh_labels;
    $self->refresh_geoconfig;
    return $self->dirty;
}

has 'dirty' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

sub _refresh_data {
    my $self = shift;
    my ($name, $file) = @_;
    my $current = $self->{$name} || {};

    my $filepath = $self->config_path . '/' . $file;

    #warn "Loading: $filepath";

    my $mtime = (stat($filepath))[9];
    unless (defined $mtime) {
        $self->log->warn("Could not read $filepath");
        return {};
    }

    if ($mtime > ($self->last_read_timestamp->{$name} || 0)) {
        my $data = $self->_read_json_safely($filepath, $current);
        $self->last_read_timestamp->{$name} = $mtime;
        if ($current ne $data) {
            $self->{$name} = $data;
            $self->dirty(1);
        }
        return $data;
    }
    return $current;
}

sub _read_json_safely {
    my ($self, $filename, $data) = @_;

    my $new = $self->_read_json($filename);
    if ($new) {
        return $new;
    }

    # keep old data
    return $data;
}

sub _read_json {
    my $self = shift;
    my $filename = shift;
    my $data = eval { decode_json(read_file($filename)) };
    $self->log->warn("Error reading $filename: $@") if $@;
    return $data;
}

1;
