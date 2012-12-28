use Test::More;

use_ok('GeoConfig');
ok(my $g = GeoConfig->new(config_path => 't/config-test'), "new");
ok(my $d = $g->dns, 'dnsconfig');

is_deeply($d->_setup_geo_rules("foo", {}), {}, "geo rules empty");
is_deeply(
    $d->_setup_geo_rules("foo", {'any1' => ''}),
    {'foo' => {a => [['10.1.1.1']]}},
    "geo rules simple"
);
is_deeply(
    $d->_setup_geo_rules("foo", {'any1' => '', 'sin1' => '10.20.1.10'}),
    {   'foo'      => {a => [['10.1.1.1']]},
        'foo.asia' => {a => [['10.20.1.10']]}
    },
    "geo rules override"
);

ok($d->setup_groups, 'setup groups');
is_deeply($d->dns->{data}->{"_edge1-global.asia"}, {a => [["10.20.1.1"]]}, "group got setup");

ok($d->setup_labels, 'setup labels');
is_deeply($d->dns->{data}->{"zone1.example"}, {alias => '_edge1-global'}, "label alias");
is_deeply($d->dns->{data}->{"zone2.example.asia"}, {a => [["10.20.1.10"]]}, "label override");

$g->pops->{sin1} = "10.20.1.101";
ok($d->setup_data, 'setup data again');
is_deeply($d->dns->{data}->{"_edge1-global.asia"}, {a => [["10.20.1.101"]]}, "new sin1 IP");


done_testing();
