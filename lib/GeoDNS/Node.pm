package GeoDNS::Node;
use Moose;

has name => (
    isa      => 'Str',
    is       => 'ro',
    required => 1,
);

has active => (
    isa     => 'Bool',
    is      => 'rw',
    default => 1,
);

has ip => (
    isa      => 'Str',
    is       => 'rw',
    required => 1,
);

1;
