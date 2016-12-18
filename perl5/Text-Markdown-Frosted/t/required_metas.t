# Test required-meta logic for Text::Markdown::Frosted::Result.
# -------------------------------------------------------------

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
    test_required_keys();
}

sub test_required_keys {
    
    
}

sub test_required_date {
    
    
}

sub test_required_string {
    
    
}

# other types?
