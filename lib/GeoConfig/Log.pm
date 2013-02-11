package GeoConfig::Log;
use v5.12.0;
use Moose::Role;
use GeoDNS::Log;

has 'log' => (
    isa     => 'Mojo::Log',
    is      => 'ro',
    default => sub { GeoDNS::Log->singleton }
);

1;
