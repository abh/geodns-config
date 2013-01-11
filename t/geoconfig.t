use Test::More;

use_ok('GeoConfig');
ok(my $g = GeoConfig->new(config_path => 't/config-test'), "new");
ok(my $pops = $g->pops, 'get pops');
ok($g->last_read_timestamp->{pops}, 'last_read is not 0');
is($pops->{"edge1.any"}, '10.1.1.1', 'any1 pop');

ok($g->refresh, 'refresh data');
is($g->pops->{"edge1.any"}, '10.1.1.1', 'any1 pop still there');

is_deeply($g->pop_geo('flex1.sin'), ['asia'], 'simple pop_geo(flex1.sin)');
is_deeply($g->pop_geo('edge1.ams'), ['nl','europe'], 'wildcard pop_geo(edge1.ams)');

done_testing();
