# Text::Markdown::Frosted::Result - Result class for Text::Markdown::Frosted.
# -------------------------------

package Text::Markdown::Frosted::Result;

use strict;
use warnings;

our $VERSION = '0.01';

use base qw(Class::Accessor::Fast);

__PACKAGE__->mk_ro_accessors(qw(file meta html input tree));

1;

__END__

=head1 NAME

Text::Markdown::Frosted::Result - Result class for Text::Markdown::Frosted.

=head1 VERSION

Version 0.01

=head1 SYNOPSIS

    use Text::Markdown::Frosted; # don't use this module directly.

=head1 DESCRIPTION

This class abstracts the result of markdown to HTML conversion performed by
C<Text::Markdown::Frosted>.  See that module for further information.

=head1 SUBROUTINES/METHODS

=head2 new ( \%params )

Create and return a new result object, with the read-only accessor values
provided in the parameter hashref.

The accessors are listed below.

=head2 file

Returns the file (path) from which the Markdown input data was read, if
applicable.  Returns C<undef> if the input data was not read from a file.

=head2 meta

Returns the metadata hashref extracted from the input data, as described in
C<Text::Markdown::Frosted>.

An empty hashref (C<{}>) is returned if the metadata section was empty or did
not exist.

=head2 html

Returns the HTML derived from the input data.

An empty string is returned if HTML document was empty, e.g. if it contained
only the meta section, contained only whitespace, or was empty.

=head2 input

Returns the input data, i.e. the raw Markdown text; or C<undef> if none was
provided.

=head2 tree

Returns the C<HTML::TreeBuilder> object used internally, if it was retained;
or C<undef> if it was not.

This can be very useful for further manipulations of the HTML.

=head1 DEPENDENCIES

=head2 C<Class::Accessor::Fast>

=head1 DIAGNOSTICS

...use the debugger as needed.

=head1 CONFIGURATION AND ENVIRONMENT

This should be entirely cross-platform and independent of configuration.
Please file a C<BUG> if you notice something.

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
