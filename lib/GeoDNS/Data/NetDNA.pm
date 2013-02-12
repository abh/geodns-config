package GeoDNS::Data::NetDNA;
use Moose::Role;

has 'last_check' => (
    isa     => 'Int',
    is      => 'rw',
    default => 0,
);

has 'api' => (
    isa      => 'NetDNA::API',
    is       => 'ro',
    required => 1,
);

has 'interval' => (
    isa     => 'Int',
    is      => 'rw',
    default => 120,
);

has 'ready' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

has 'dirty' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

sub check {
    my $self = shift;
    if ($self->last_check + $self->interval > time) {
        return 0;
    }
    else {
        return $self->update;
    }
}



1;
