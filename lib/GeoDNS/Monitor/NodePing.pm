package GeoDNS::Monitor::NodePing;
use v5.12.0;
use Moose;
extends 'GeoDNS::Monitor';
with 'GeoDNS::Log::Role';
use Data::Dump qw(pp);
use JSON qw(decode_json);

use Mojo::UserAgent;
use Mojo::URL;

has api_key => (
    isa      => 'Str',
    is       => 'rw',
    required => 1,
);

has base_url => (
    isa     => 'Str',
    is      => 'rw',
    default => 'https://api.nodeping.com/api/1/',
);

has account => (
    isa     => 'Str',
    is      => 'rw',
    default => ''
);

has ua => (
    isa     => 'Mojo::UserAgent',
    is      => 'rw',
    default => sub { Mojo::UserAgent->new }
);

has 'config' => (
    isa      => 'GeoConfig',
    is       => 'ro',
    required => 1,
    weak_ref => 1
);

has 'dirty' => (
    isa     => 'Bool',
    is      => 'rw',
    default => 1,
);

has 'checks' => (
    isa     => 'HashRef[HashRef]',
    is      => 'rw',
    default => sub { {} }
);

has fqdn_suffix => (
    isa => 'Str',
    is => 'ro',
    default => '',
);

sub url {
    my ($self, $method, $params) = @_;
    my $url = Mojo::URL->new($self->base_url . $method);
    $url->userinfo($self->api_key);
    if ($params) {
        $url->query->param(%$params);
    }
    if ($self->account) {
        $url->query->param(customerid => $self->account);
    }
    say "URL: ", $url->to_string;
    return $url;
}

sub get {
    my ($self, $method, $params, $cb) = @_;
    if (ref $params eq 'CODE') {
        $cb     = $params;
        $params = undef;
    }
    $params ||= {};
    $self->ua->get(
        $self->url($method, $params),
        sub {
            my ($ua, $tx) = @_;

            #pp($tx->res);
            #pp($tx->res->content->headers);
            say "BODY: ", $tx->res->body;
            my $data = $tx->res->json;
            $cb->($data) if $cb;
        }
    );

}

sub post {
    my ($self, $method, $params, $cb) = @_;
    my $url = $self->url($method);
    return $self->ua->post($url, {}, json => $params, $cb);
}


sub start {
    my $self = shift;

    my $delay = Mojo::IOLoop->delay(
        sub {
            my $delay = shift;
            $self->load_checks($delay->begin);
        },
        sub {
            my $delay = shift;
            $self->load_outages;
            my $end = $delay->begin;
            Mojo::IOLoop->timer(4 => sub {
                warn "going to sync checks";
                $self->sync_checks($end);
                $self->load_outages;
            });
        },
        sub {
            Mojo::IOLoop->recurring(
                50 => sub {
                    my $loop = shift;
                    warn "loading outages";
                    $self->load_outages;
                }
            );
        },
        sub {
            Mojo::IOLoop->recurring(
                600 => sub {
                    my $loop = shift;
                    warn "syncing checks";
                    $self->sync_checks($loop->delay->begin);
                }
            );
        }
    );
}

sub sync_checks {
    my ($self, $cb) = @_;
    return if $self->dirty;
    my $nodes = $self->config->nodes;

    unless ($nodes->ready) {
        $cb->() if $cb;
        return;
    }

    my $checks;

    say "NODES: ", pp($nodes->nodes);

    my $delay = Mojo::IOLoop->delay;

    while (my ($name, $node) = each %{$nodes->nodes}) {
        next unless $node->active;
        # $self->log->info("checking monitoring for $name");
        next if $self->is_monitored($name);
        $self->log->info("$name isn't monitored, yet");

        my $server = $name;
        if ($name !~ m/\.$/ and $self->fqdn_suffix) {
            $server = $name . "." . $self->fqdn_suffix;
        }

        my $ip = $node->ip;

        $self->log->info("Adding monitoring for $name/$ip");

        #$cb->() if $cb;
        #return;


        $self->post("checks", {
           label         => "$name",
            interval      => 1,
            enabled       => 'active',
            runlocations  => 'wlw',
            public        => 'false',
            threshold     => 2,
            sens          => 5,
            type         => 'HTTP',
            target       => "http://$ip/",
        }, $delay->begin);

        $self->checks->{$name} = {};

    }
    $cb->() if $cb;
}

sub load_outages {
    my ($self, $cb) = @_;
      $self->get(
        'results?action=current',
        sub {
            my $results = shift;
            say "RESULTS: ", pp($results);
            for my $check (values %{$self->checks}) {
                next unless $check->{_id};
                $check->{result} = $results->{$check->{_id}} || {};
            }
            $cb->() if $cb;
        }
    );
}

sub load_checks {
    my ($self, $cb) = @_;

    $self->get(
        'checks',
        sub {
            my $checks = shift;
            say "CHECKS: ", pp($checks);
            $checks = {} if ref $checks eq 'ARRAY';
            for my $check (values %$checks) {
                $self->checks->{$check->{label}} = $check;
            }
            $self->ready(1);
            $self->dirty(0);
            $self->sync_checks();
            $cb->() if $cb;
        }
    );
}

sub status {
    my ($self, $node) = @_;
    my $check = $self->check($node) or return { status => "none" };
    my $result = $check->{result}   or return { status => "up" };

    return {
        status => ($result->{type}),
        message => ($result->{reason}
            ? ($result->{reason}->{message} || $result->{reason}->{code})
            : "")
    };
}

sub is_monitored {
    my ($self, $node) = @_;
    return 1 if $self->check($node);
    return 0;
}

sub check {
    my ($self, $node) = @_;
    return $self->checks->{$node};
}

1;
