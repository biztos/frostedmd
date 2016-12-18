# Text::Markdown::Frosted::TOML::Parser - modified TOML parser.
# -------------------------------------

package Text::Markdown::Frosted::TOML::Parser;

use strict;
use warnings;

use base qw(TOML::Parser);

our $VERSION = '0.01';

sub _tokenizer_class {
    return 'Text::Markdown::Frosted::TOML::Tokenizer';
}

1;

__END__

=head1 NAME

Text::Markdown::Frosted::TOML::Parser - modified TOML parser.

=head1 VERSION

Version 0.01

=head1 SYNOPSIS

    use Text::Markdown::Frosted; # don't use this module directly.

=head1 DESCRIPTION

This is a subclass of C<TOML::Parser>, enabling the use of a custom tokenizer
which in turn allows the more flexible date specification we expect to find
in TOML 1.0.

L<https://github.com/toml-lang/toml/pull/297>

As this class overrides a "private" method, it should be considered
I<fragile> at best.

=head1 EXPECTED DEPRECATION

We expect this module to become obsolete when TOML 1.0 is released and the
C<TOML> module is brought into sync with it.

=head1 SUBROUTINES/METHODS

=head2 _tokenizer_class

Overridden to always return C<Text::Markdown::Frosted::TOML::Tokenizer>.

=head1 DEPENDENCIES

=head2 C<TOML::Parser>

=head2 C<TOML::Parser::Tokenizer>

=head2 C<Text::Markdown::Frosted::TOML::Tokenizer>

=head1 DIAGNOSTICS

C<TOML::Parser> supports the C<TOML_PARSER_TOKENIZER_DEBUG> environment
variable, which is very useful.

=head1 CONFIGURATION AND ENVIRONMENT

As per C<TOML::Parser>; presumably nothing special.

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
