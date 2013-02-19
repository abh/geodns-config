use Test::More;

use_ok('GeoDNS::Labels::File');
ok(my $labels = GeoDNS::Labels::File->new(file => 't/config-test/labels.json'), "new");
isa_ok($labels, 'GeoDNS::Labels');
ok($labels->update, "update");
ok($labels->last_read_timestamp, "last_read_timestamp");

done_testing();
