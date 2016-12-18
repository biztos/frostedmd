# Text::Markdown::Frosted::TOML::Tokenizer - modified TOML tokenizer.
# ----------------------------------------

package Text::Markdown::Frosted::TOML::Tokenizer;

use strict;
use warnings;

use Carp qw(confess);

use base qw(TOML::Parser::Tokenizer);

our $VERSION = '0.01';

sub grammar_regexp {

    my $class = shift;

    # Only mess with the specific datetime regexp.
    my $data = $class->SUPER::grammar_regexp();
    confess 'unexpected grammar_regexp format in TOML::Parser::Tokenizer'
        unless $data->{value}->{datetime};

    $data->{value}->{datetime} = qr{
                (
                    [0-9]{4}- # YYYY
                    [0-9]{2}- # MM
                    [0-9]{2}  # DD
                    T
                    [0-9]{2}: # HH
                    [0-9]{2}: # MI
                    [0-9]{2}  # SS
                    (?:[.][0-9]+)? # seconds can be a float
                    (?:
                        Z
                        |
                        -[0-9]{2}:[0-9]{2} # TZ
                    )
                )
            }x;

    return $data;

}

1;

__END__

=head1 NAME

Text::Markdown::Frosted::TOML::Tokenizer - modified TOML tokenizer.

=head1 VERSION

Version 0.01

=head1 SYNOPSIS

    use Text::Markdown::Frosted; # don't use this module directly.

=head1 DESCRIPTION

This is a subclass of C<TOML::Parser::Tokenizer>, enabling the more flexible
date types we expect to find in TOML 1.0.

L<https://github.com/toml-lang/toml/pull/297>

=head1 SUBROUTINES/METHODS

=head2 grammar_regexp

Overridden to return a hashref that includes a more liberal datetime regex.

Nothing else is changed in the hashref.

=head1 DEPENDENCIES

=head2 C<TOML::Parser::Tokenizer>

=head1 DIAGNOSTICS

C<TOML::Tokenizer> supports the C<TOML_PARSER_TOKENIZER_DEBUG> environment
variable, which is very useful.

=head1 CONFIGURATION AND ENVIRONMENT

As per C<TOML::Parser::Tokenizer>; presumably nothing special.

=head1 INCOMPATIBILITIES

None; please file a C<BUG> if you find any.

=head1 AUTHOR

Kevin Frost, C<< <biztos at mac.com> >>

=head1 BUGS AND LIMITATIONS

Please file bugs and feature requests in the GitHub issue tracker:

L<https://github.com/biztos/frosted-markdown/issues>

That helps make sure bugs are fixed, and useful features added, across the
Frosted Markdown implementations.

If you prefer, you may also report bugs via the CPAN interface by sending
e-mail to C<bug-text-markdown-frosted at rt.cpan.org>, or through the web
interface at
L<http://rt.cpan.org/NoAuth/ReportBug.html?Queue=Text-Markdown-Frosted>.

=head1 LICENSE AND COPYRIGHT

Copyright 2015 Kevin Frost.

This program is free software; you can redistribute it and/or modify it
under the terms of the the Artistic License (2.0). You may obtain a
copy of the full license at:

L<http://www.perlfoundation.org/artistic_license_2_0>

For details see the LICENSE section in C<Text::Markdown::Frosted>.

=cut
