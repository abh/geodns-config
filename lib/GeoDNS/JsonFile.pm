package GeoDNS::JsonFile;
use Moose::Role;
with 'GeoDNS::JsonFiles';

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
    my $self = shift;
    if ($self->_refresh_data('data', $self->file)) {
        $self->ready(1);
        return 1;
    }
    return 0;
}

sub all {
    my $self = shift;
    return \%{$self->{data}};
}

1;
