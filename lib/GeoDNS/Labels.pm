package GeoDNS::Labels;
use Moose;

sub all {
    my $self = shift;
    return \%{$self->{data}};
}

1;
