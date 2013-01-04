package GeoDNS::Monitor::Panopta;
use v5.12.0;
use Moose;
use Data::Dump qw(pp);
use JSON qw(decode_json);

has mojo => (
    required => 1,
    is       => 'ro',
);

has username => (
    isa      => 'Str',
    is       => 'ro',
    required => 1,
);

has password => (
    isa      => 'Str',
    is       => 'ro',
    required => 1,
);

has ua => (
    isa        => 'Mojo::UserAgent',
    is         => 'ro',
    lazy_build => 1
);

sub _build_ua {
    my $ua = shift->mojo->ua();
    $ua->inactivity_timeout(20);
    return $ua;
}

has servers => (
    isa => 'HashRef',
    is  => 'rw',
    default => sub { {} }
);

has outages_list => (
    isa => 'HashRef',
    is  => 'rw',
    default => sub { {} }
);

sub _r {
    my ($self, $section, $operation, $args, $cb) = @_;
    if (ref $args eq 'CODE') {
        $cb   = $args;
        $args = undef;
    }
    $args ||= {};

    my $url ="https://api.panopta.com/$section/$operation";

    #$url = "http://www.yellowbot.com/api/test/echo";

    $self->mojo->log->debug("URL: ", $url);
    
    my $r = $self->ua->post_form(
        $url,
        {   username => $self->username,
            password => $self->password,
            %$args
        },
        sub {
            my ($ua, $tx) = @_;
            #pp($tx->res);
            #pp($tx->res->content->headers);
            say "BODY: ",$tx->res->body;
            my $data = decode_json($tx->res->body);
            #say "DATA: ", pp($data);
            $cb->($data) if $cb;
        }
    );

}

sub start {
    my $self= shift;
    $self->mojo->log->info("Starting panopta ...");
    $self->load_servers;
    Mojo::IOLoop->recurring(
        120 => sub {
            $self->load_outages;
        }
    );
    Mojo::IOLoop->recurring(
        1800 => sub {
            $self->load_servers;
        }
    );
}

sub load_servers {
    my $self = shift;
    $self->_r(
        "config",
        "listServers",
        {},
        sub {
            my $data = shift;
            $self->servers({map { $_->{server_id} => $_ } @{$data->{servers}}});
            $self->load_outages;
        }
    );
}

sub load_outages {
    my $self = shift;
    $self->_r(
        "status",
        "getCurrentOutages",
        {},
        sub {
            my $data = shift;
            $self->outages_list({map { $_->{server_id} => $_ } @{$data->{outages}}});
        }
    );
}

sub ip {

}

sub ips {
    my $self = shift;
    my $servers = $self->servers;
}

sub outages {
    my $self = shift;
    my $list = $self->outages_list;
    my $servers = $self->servers;
    my %outages;
    for my $server_id (sort keys %$list) {
        $outages{$servers->{$server_id}->{last_known_ip}} = $list->{$server_id};
    }
    return \%outages;
}

1;
