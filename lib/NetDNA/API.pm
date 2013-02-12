package NetDNA::API;
use Moose;
use Net::OAuth;
use Mojo::UserAgent;

$Net::OAuth::PROTOCOL_VERSION = Net::OAuth::PROTOCOL_VERSION_1_0A;

has [ 'alias', 'key', 'secret' ] => (
   isa => 'Str',
   is  => 'rw',
   required => 1,
);

has 'host' => (
   isa => 'Str',
   is  => 'rw',
   default => 'rws.netdna.com',
);

has 'ua' => (
   isa => 'Mojo::UserAgent',
   is  => 'rw',
   default => sub { Mojo::UserAgent->new }, 
);

#get("flexnodes.json");
#get("nodes.json");

sub url {
    my ($self, $method) = @_;
    
    my $address = join "/", 'https:/', $self->host, $self->alias, $method;

    # Create request
    my $request = Net::OAuth->request("request token")->new(
        consumer_key     => $self->key,
        consumer_secret  => $self->secret,
        request_url      => $address,
        request_method   => 'GET',
        signature_method => 'HMAC-SHA1',
        timestamp        => time,
        nonce            => '',
        callback         => '',
    );

    $request->sign;
    return $request->to_url;
}

sub get {
    my ($self, $method) = @_;
    my $ua = Mojo::UserAgent->new();
    return $ua->get($self->url($method))->res;
}

1;
