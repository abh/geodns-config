package GeoDNS::Labels::File;
use Moose;
extends 'GeoDNS::Labels';
with 'GeoDNS::JsonFile', 'GeoConfig::Log';

has 'name' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'labels',
);

1;
