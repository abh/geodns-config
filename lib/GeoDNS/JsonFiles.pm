package GeoDNS::JsonFiles;
use Moose::Role;
use JSON qw(decode_json encode_json);
use File::Slurp qw(read_file write_file);

# TODO:
#   When GeoConfig isn't reading things directly anymore, fold
#   this code into JsonFile (without the "multifile support")

use namespace::clean;

has 'dirty' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

has 'last_read_timestamp' => (
    isa     => 'HashRef',
    is      => 'rw',
    default => sub { return {} }
);

sub _refresh_data {
    my $self = shift;
    my ($name, $file) = @_;
    my $current = $self->{$name} || {};

    my $filepath = $self->can('config_path') ? $self->config_path . '/' . $file : $file;

    my $mtime = (stat($filepath))[9];
    unless (defined $mtime) {
        $self->log->warn("Could not read $filepath");
        return 0;
    }

    if ($mtime > ($self->last_read_timestamp->{$name} || 0)) {
        my $data = $self->_read_json_safely($filepath, $current);
        $self->last_read_timestamp->{$name} = $mtime;
        if ($current ne $data) {
            $self->log->info("Loaded $filepath");
            $self->{$name} = $data;
            $self->dirty(1);
        }
        return 1;
    }
    return 0;
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
