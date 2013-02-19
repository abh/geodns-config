package GeoDNS::JsonFile;
use Moose::Role;
use JSON qw(decode_json encode_json);
use File::Slurp qw(read_file write_file);

use namespace::clean;

has 'dirty' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

has 'last_read_timestamp' => (
    isa     => 'Int',
    is      => 'rw',
    default => 0,
);

has 'file' => (
    isa     => 'Str',
    is      => 'rw',
    lazy    => 1,
    default => sub { shift->name . '.json' },
);

has 'ready' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

sub check {
    my $self = shift;
    return $self->update;
}

sub update {
    my ($self, $cb) = @_;
    if ($self->_refresh_data($self->file, $cb)) {
        unless ($cb) {
            $self->ready(1);
            return 1;
        }
    }
    return 0;
}

sub _refresh_data {
    my $self = shift;
    my ($file, $cb) = @_;

    my $current = $self->{_json_data} || {};

    my $filepath = $self->can('config_path') ? $self->config_path . '/' . $file : $file;

    my $mtime = (stat($filepath))[9];
    unless (defined $mtime) {
        $self->log->warn("Could not read $filepath");
        return 0;
    }

    if ($mtime > ($self->last_read_timestamp || 0)) {
        my $data = $self->_read_json_safely($filepath, $current);
        $self->last_read_timestamp($mtime);
        if ($current ne $data) {
            $self->log->info("Loaded $filepath");
            $self->dirty(1);
            $self->{_json_data} = $data;
            if ($cb) {
                $cb->($data);
            }
            else {
                $self->{data} = $data
            }
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
