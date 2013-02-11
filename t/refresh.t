use Test::More;
use File::Copy qw(copy);
use File::Slurp qw(write_file);

use_ok('GeoConfig');
ok(my $g = GeoConfig->new(domain_name => 'example.net', config_path => 't/config-test'), "new");
ok($g->nodes);
ok($g->refresh,         'refresh');
ok($g->dns->setup_data, 'setup data');
ok($g->dns->dns->{data}->{"_edge1-global"}, "has data setup");
Data::Dump::pp($g->dns->dns);
is($g->nodes->node_ip("edge1.any"), '10.1.1.1', 'any1 pop');

my $labels_file = "t/config-test/labels.json";

copy($labels_file, "${labels_file}.bak")
  or die "Could not copy labels.json: $!";

write_file($labels_file, '{
  "zone1.example":{"group":"edge1-global"},
  "zone2.example":{"any1":"10.1.1.10","sin1":"10.20.1.10"},
  "another":{"group":"edge1-global"}
}');

ok($g->refresh, 'refresh');

ok($g->dns->setup_data, 'setup data');
ok($g->dns->dns->{data}->{"_edge1-global"}, "has data setup");
ok($g->dns->dns->{data}->{"another"}, "has new data setup");

copy("${labels_file}.bak", $labels_file)
  or die "Could not copy labels.json: $!";

done_testing();
