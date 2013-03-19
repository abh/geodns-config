package GeoDNS::Monitor;
use Moose;

has 'ready' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 0,
);

1;
