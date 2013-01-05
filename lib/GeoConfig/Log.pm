package GeoConfig::Log;
use v5.12.0;
use Moose;
use Mojo::Log;

has 'log_level' => (
  isa => 'Str',
  is => 'ro',
  default  => 'debug'
);

has 'log_path' => (
  isa => 'Str',
  is => 'ro',
  default  => 'STDERR'
);

has 'log' => (
 isa => 'Mojo::Log',
 is  => 'ro', 
 default => sub { Mojo::Log->new(level => $_[0]->log_level, path => $_[0]->log_path) } 
); 

1;
