# (Minimal) tests for Text::Markdown::Frosted::Result.
# ----------------------------------------------------
# Interestingly, this is one of those classes that's so simple, in most
# strongly-typed languages it would be redundant to test it at all.

use strict;
use warnings;

use Test::Exception;

use Test::More tests => 50;

my $CLASS;

BEGIN {
    $CLASS = 'Text::Markdown::Frosted::Result';
    use_ok($CLASS) or die "I give up";
}

MAIN: {
    test_construction_with_accessors();
    test_ro_accessors();
}

sub test_construction_with_accessors {

    note 'test_construction_with_accessors';

    my $res;

    # Empty:
    lives_ok { $res = $CLASS->new() } 'no params required';
    isa_ok( $res, $CLASS, 'result' );

    # Full:
    my $possible_params = {
        file  => 'some/file.md',
        meta  => { ima => 'meta' },
        html  => '<h1>here</h1>',
        input => '# here',
        tree  => 'yep, a tree',
    };
    lives_ok { $res = $CLASS->new() } 'full params accepted';
    isa_ok( $res, $CLASS, 'result' );

    # With/without each:
    for my $p ( keys %$possible_params ) {
        lives_ok { $res = $CLASS->new( { $p => $$ } ) } "only $p";
        isa_ok( $res, $CLASS, 'result' );
        is( $res->$p, $$, "$p sticks" );

        my $params = { %{$possible_params} };
        delete $params->{$p};
        lives_ok { $res = $CLASS->new($params) } "all but $p";
        isa_ok( $res, $CLASS, 'result' );
        is( $res->$p, undef, "$p undef" );

    }
}

sub test_ro_accessors {

    note 'test_ro_accessors';

    my @accessors = qw(file meta html input tree);

    for my $acc (@accessors) {
        my $res;
        lives_ok { $res = $CLASS->new() } 'no params required';
        isa_ok( $res, $CLASS, 'result' );
        throws_ok { $res->$acc('anything') }
        qr/cannot alter the value of '$acc'/, "$acc is r/o";

    }
}
