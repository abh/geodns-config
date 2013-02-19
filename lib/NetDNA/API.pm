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
    my ($self, $method, $cb) = @_;

    $self->ua->get(
        $self->url($method),
        sub {
            my ($ua, $tx) = @_;
            if ($cb) {
                $cb->($tx->res);
            }
            else {
                # todo make sure this is useful and log the response appropriately if it is
                $self->log->info("$method called without callback");
            }
        }
    );
}

1;
