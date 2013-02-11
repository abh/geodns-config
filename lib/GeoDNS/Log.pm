package GeoDNS::Log;
use strict;
use Mojo::Log;

my $log = Mojo::Log->new;

sub singleton {
    return $log; 
}

1;
