use Test::More;

use_ok('GeoDNS::Monitor::Manual');

ok(my $monitor = GeoDNS::Monitor::Manual->new());
ok($monitor->set_outage('10.0.0.1', 1));
ok($monitor->outage('10.0.0.1'));
ok($monitor->set_outage('10.0.0.1', 0)==0);
is($monitor->outage('10.0.0.1'),0);

done_testing();