use Test::More;

use_ok('GeoConfig');
ok(my $g = GeoConfig->new(config_path => 't/config-test'), "new");
ok(my $pops = $g->pops, 'get pops');
ok($g->last_read_timestamp->{pops}, 'last_read is not 0');
is($pops->{any1}, '10.1.1.1', 'any1 pop');

ok($g->refresh, 'refresh data');
is($g->pops->{any1}, '10.1.1.1', 'any1 pop still there');

done_testing();
