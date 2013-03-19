package GeoDNS::Monitor::Manual;
use v5.12.0;
use Moose;
extends 'GeoDNS::Monitor';

has 'ips' => (
	traits  => ['Hash'],
	isa => 'HashRef[Str]',
	is  => 'rw',
	default => sub { {} },
	handles => {
		'set_outage' => 'set',
		'outage' => 'get',
	}
);

1;