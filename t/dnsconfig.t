use Test::More;

use Data::Dump qw(pp);

use_ok('GeoConfig');
use_ok('GeoDNS::Monitor::Manual');
ok( my $g = GeoConfig->new(
        domain_name => 'example.net',
        config_path => 't/config-test',
    ),
    "new"
);
ok(my $d = $g->dns, 'dnsconfig');

# test outage
ok($d->config->monitor->set_outage("10.20.1.1",  1), "setting outage");
ok($d->setup_data, 'setup labels');
is_deeply($d->dns->{data}->{"zone2.example.asia"}, undef, "outage disabled asia");

# outage over
ok($d->config->monitor->set_outage("10.20.1.1", 0) == 0, "clearing outage");
ok($d->setup_data, 'setup labels');
is_deeply($d->dns->{data}->{"zone2.example.asia"}, {a => [["10.20.1.10"]]}, "asia is back");

# test changing the POP ip
ok($g->nodes->set_ip("flex1.sin", "10.20.1.101"), "update sin1 IP");
ok($g->nodes->node_ip("flex1.sin", "10.20.1.101"), "sin1 IP changed");
ok($d->setup_data, 'setup data again');
is_deeply($d->dns->{data}->{"_edge1-global.asia"}, {a => [["10.20.1.101"]]}, "new sin1 IP");
diag(pp($d->dns));

done_testing();
