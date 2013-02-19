use Test::More;

use_ok('GeoDNS::Nodes::File');
ok(my $nodes = GeoDNS::Nodes::File->new(name => 't/config-test/pops'), "new");
ok($nodes->update, 'update');
is($nodes->nodes->{'edge1.any'} && $nodes->nodes->{'edge1.any'}->ip, '10.1.1.1', 'edge1.any ip');
isa_ok($nodes, 'GeoDNS::Nodes');

done_testing();
