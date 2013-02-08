use Test::More;

use_ok('GeoConfig');
use_ok('GeoDNS::Monitor::Manual');
ok( my $g = GeoConfig->new(
        config_path => 't/config-test',
        monitor     => GeoDNS::Monitor::Manual->new
    ),
    "new"
);
ok(my $d = $g->dns, 'dnsconfig');

is_deeply($d->_setup_geo_rules("foo", {}), {}, "geo rules empty");
is_deeply(
    $d->_setup_geo_rules("foo", {'edge1.any' => ''}),
    {'foo' => {a => [['10.1.1.1']]}},
    "geo rules simple"
);
is_deeply(
    $d->_setup_geo_rules("foo", {'edge1.any' => '', 'flex1.sin' => '10.20.1.10'}),
    {   'foo'      => {a => [['10.1.1.1']]},
        'foo.asia' => {a => [['10.20.1.10']]}
    },
    "geo rules override"
);

done_testing();
