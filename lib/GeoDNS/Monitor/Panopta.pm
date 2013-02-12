package GeoDNS::Monitor::Panopta;
use v5.12.0;
use Moose;
extends 'GeoDNS::Monitor', 'GeoDNS::Log';
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
    isa     => 'HashRef',
    is      => 'rw',
    default => sub { {} }
);

has outages_list => (
    isa     => 'HashRef',
    is      => 'rw',
    default => sub { {} }
);

sub _r {
    my ($self, $section, $operation, $args, $cb) = @_;
    if (ref $args eq 'CODE') {
        $cb   = $args;
        $args = undef;
    }
    $args ||= {};

    my $url = "https://api.panopta.com/$section/$operation";

    #$self->mojo->log->debug("Panopta URL: ", $url);

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
            #say "BODY: ", $tx->res->body;
            my $data = decode_json($tx->res->body);

            #say "DATA: ", pp($data);
            $cb->($data) if $cb;
        }
    );

}

sub start {
    my $self = shift;
    my $delay = Mojo::IOLoop->delay(
        sub {
            my $delay = shift;
            $self->load_servers($delay->begin);
        },
        sub {
            my $delay = shift;
            $self->load_outages($delay->begin);
        },
        sub {
            Mojo::IOLoop->recurring(
                60 => sub {
                    $self->load_outages;
                }
            );

            Mojo::IOLoop->recurring(
                1800 => sub {
                    $self->load_servers;
                }
            );

        }
    );

}

sub load_servers {
    my ($self, $cb) = @_;
    $self->mojo->log->info("Refreshing Panopta server list");
    $self->_r(
        "config",
        "listServers",
        {},
        sub {
            my $data = shift;
            $self->servers({map { $_->{server_id} => $_ } @{$data->{servers}}});
            if ($cb) { $cb->() }
        }
    );
}

sub load_outages {
    my ($self, $cb) = @_;
    $self->mojo->log->info("Refreshing Panopta outages");
    $self->_r(
        "status",
        "getCurrentOutages",
        {},
        sub {
            my $data = shift;
            $self->outages_list({map { $_->{server_id} => $_ } @{$data->{outages}}});
            if ($cb) { $cb->() }
        }
    );
}

sub check {
    my ($self, $id) = @_;
    my $servers = $self->servers;
    return $servers->{$id} || {};
}

sub ips {
    my $self    = shift;
    my $servers = $self->servers;
    my $outages = $self->outages;
    my %ips;
    for my $server_id (sort keys %$servers) {
        my $ip = $servers->{$server_id}->{last_known_ip};
        if (!$ips{$ip}) {
            $ips{$ip} = [];
        }
        push @{$ips{$ip}},
          { id    => $server_id,
            name  => ($servers->{$server_id}->{name} || ''),
            alert => ($outages->{$server_id} ? 1 : 0)
          };
    }
    return \%ips;
}

sub outages {
    my $self    = shift;
    my $list    = $self->outages_list;
    my $servers = $self->servers;
    my %outages;
    for my $server_id (sort keys %$list) {
        $outages{$servers->{$server_id}->{last_known_ip}} = $list->{$server_id};
    }
    return \%outages;
}

sub outage {
    my $self = shift;
    my $ip = shift;
    return $self->outages->{$ip} ? 1 : 0;
}

1;
