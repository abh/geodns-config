use Test::More;
use Data::Dump qw(pp);

use_ok('GeoConfig');
ok(my $g = GeoConfig->new(domain_name => 'example.net', config_path => 't/config-test'), "new");
ok($g->refresh, 'refresh data');
is($g->nodes->update, 0, 'update nodes, not changed');
isnt($g->nodes->last_read_timestamp, 0, 'last_read is set');

is($g->nodes->node_ip("edge1.any"), '10.1.1.1', 'any1 pop');

# TODO: need to simulate the data changing
ok($g->nodes->update==0, 'refresh data');
is($g->nodes->node_ip("edge1.any"), '10.1.1.1', 'any1 pop still there');

is_deeply($g->pop_geo('flex1.sin'), ['asia'], 'simple pop_geo(flex1.sin)');
is_deeply($g->pop_geo('edge1.ams'), ['nl','europe'], 'wildcard pop_geo(edge1.ams)');

done_testing();
