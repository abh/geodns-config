package GeoDNS::Labels::File;
use Moose;
extends 'GeoDNS::Labels';
with 'GeoDNS::JsonFile', 'GeoDNS::Log::Role';

has 'name' => (
    isa     => 'Str',
    is      => 'ro',
    default => 'labels',
);

1;
