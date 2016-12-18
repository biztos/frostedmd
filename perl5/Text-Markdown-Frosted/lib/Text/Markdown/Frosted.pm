package Text::Markdown::Frosted;

use strict;
use warnings;

use Carp qw(croak confess);
use IO::File;
use Text::Markdown;
use HTML::TreeBuilder;

use Text::Markdown::Frosted::Result;
use Text::Markdown::Frosted::TOML::Parser;

use base qw(Class::Accessor::Fast);

__PACKAGE__->mk_ro_accessors(qw(md_opts keep_tree required meta_at_bottom));
__PACKAGE__->mk_accessors(qw(error));

sub new {

    my ( $class, $params ) = @_;
    $params = { %{ $params // {} } };

    # All non-local opts are passed through to Text::Markdown.
    if ( !$params->{md_opts} ) {
        my $is_fm = { map { $_ => 1 } qw(keep_tree required meta_at_bottom) };
        my $md_opts = {};
        for my $key ( keys %{$params} ) {
            if ( !$is_fm->{$key} ) {
                $md_opts->{$key} = delete $params->{$key};
            }
        }
        $params->{md_opts} = $md_opts;
    }

    if ( my $required = $params->{required} ) {
        croak 'required must be a hashref' if ref $required ne 'HASH';
        $params->{required} = _init_required($required);
    }

    return $class->SUPER::new($params);

}

sub _init_required_types {

    my $required = { %{shift} };

    for my $key(sort keys %{$required}) {
        my $val = $required->{$key};
        if ($val eq 'date' || $val eq 'datetime') {
            $required->{ref} = 'DateTime';
        }
#        elsif ($val eq 'int')
 #       if ($key eq 'date' || $key 
    }
}

sub markdown {

    my ( $data, $opts ) = @_;

    my $fm = __PACKAGE__->new($opts);

    my $res = $fm->convert($data) or croak $fm->error;
    if (wantarray) {
        return ( $res->meta, $res->html );
    }
    else {
        return $res->html;
    }
}

sub markdown_file {

    my ( $file, $opts ) = @_;

    my $fm = __PACKAGE__->new($opts);

    my $res = $fm->convert_file($file) or croak $fm->error;
    if (wantarray) {
        return ( $res->meta, $res->html );
    }
    else {
        return $res->html;
    }

}

sub convert_file {

    my ( $self, $file ) = @_;

    my $fh = IO::File->new( $file, '<' )
        or croak "Failed to open file $file for read: $!";
    my $data = do { local $/; <$fh> };
    close $fh;

    return $self->_convert( $data, $file );

}

sub convert {

    my ( $self, $data ) = @_;

    return $self->_convert($data);

}

sub _convert {

    my ( $self, $data, $file ) = @_;

    # $file is optional, but $data has to at least be defined.
    if ( !defined $data ) {
        $self->error('no data to convert');
        return;
    }

    # clear any previous error just in case
    $self->error(undef);

    my $md_html = Text::Markdown::markdown( $data, $self->md_opts );

    my $tree = HTML::TreeBuilder->new();

    $tree->parse($md_html);

    my $meta = $self->_frost_tree($tree) or return;

    my $res_params = {
        file  => $file,
        input => $data,
        meta  => $meta,
        html  => $tree->as_HTML,
    };

    if ( $self->keep_tree ) { $res_params->{tree} = $tree; }

    return Text::Markdown::Frosted::Result->new($res_params);

}

sub _frost_tree {

    my ( $self, $tree ) = @_;

    # Convert and remove the meta block if it's there.
    my $meta = {};
    if ( my $node = $self->_find_meta_node($tree) ) {
        my $parser = Text::Markdown::Frosted::TOML::Parser->new();
        my $meta = eval { $parser->parse( $node->content_list ) };
        if ( my $err = $@ ) {
            $self->error( 'TOML parse error: ' . $err );
            return;
        }
        $node->parent->detach();
    }

    # Check for compliance with expected data structures (if applicable).
    if ( $self->required ) {
        $self->_enforce_required_meta() or return;
    }

    # Convert the image nodes in place.

    return $meta;

}

sub _find_meta_node {

    my ( $self, $tree ) = @_;

    if ( !$tree || !ref $tree || $tree->tag ne 'html' ) {
        confess 'malformed or unexpected tree: top not "html"';
    }
    my @top_elems = $tree->content_list;
    my $body      = $top_elems[1];
    if ( !$body || !ref $body || $body->tag ne 'body' ) {
        confess 'malformed or unexpected tree: body not found';
    }

    # The meta block must be a direct child of the body, and it has to come
    # either at the "top" of the document or at the "bottom."  The top is
    # somewhat flexible: normally you might have an h1 before it, or even
    # more; but it must come before any non-heading element.
    #
    # The bottom is stricter: it must be the last non-whitespace thing in the
    # body (excluding whitespace).
    if ( $self->meta_at_bottom ) {
        for my $elem ( reverse $body->content_list ) {
            if ( !ref $elem ) {
                if ( $elem =~ /\S/s ) {
                    return;    # hit non-whitespace text (how?)
                }
                next;
            }
            elsif ( my $meta_node = $self->_pre_code($elem) ) {
                return $meta_node;
            }
            else {
                return;        # hit another element.
            }
        }
        return;                # no elements (odd but possible).
    }

    # Try it from the top.
    my $headings = 0;
    my $max = $self->max_headings // 1;
    for my $elem ( $body->content_list ) {
        next unless ref $elem;    # skip any stray text nodes
        if ( my $meta_node = $self->_pre_code($elem) ) {
            return $meta_node;
        }
        if ( $elem->tag =~ /^h[\d]$/i ) {
            if ( ++$headings > $max ) {
                return;           # too far down in the document.
            }
        }

    }

    return;                       # got nothin'

}

sub _pre_code {

    my ( $self, $node ) = @_;

    return unless $node && ref $node && $node->tag eq 'pre';

    my @kids = $node->content_list;
    return unless @kids == 1 && ref $kids[0] && $kids[0]->tag eq 'code';

    return $kids[0];

}

sub _enforce_required_meta {

    my ( $self, $meta ) = @_;

    my $req = $self->required;
    croak 'required must be a hashref' if ref $req ne 'HASH';

    # All required items must be present and have the correct type.
    #
    # TOML is case sensitive and thus so are we.
    #
    # NOTE: this only applies to the top level.  If you need to validate
    # data deeper than this you should do it in the caller.
    my @check_keys = sort keys %{$req};
    my $check      = { %{$req} };
    for my $key ( sort keys %{$meta} ) {
        if ( $check->{$key} ) {
            $check->{$key}->{have} = 1;

            # TODO: TYPE CHECK!
        }
    }
    my @missing = grep { !$check->{$_}->{have} } @check_keys;
    my @bad     = grep { !$check->{$_}->{ok} } @check_keys;
    my @errs;
    push @errs, sprintf( 'Missing required key: %s', $_ ) for @missing;
    push @errs, sprintf( 'Wrong format for key: %s', $_ ) for @bad;
    if (@errs) {
        $self->err( join '; ', @errs );
        return;
    }

    return 1;
}

=head1 NAME

Text::Markdown::Frosted - Markdown to HTML with Frosting on top!

=head1 VERSION

Version 0.01

=cut

our $VERSION = '0.01';

=head1 SYNOPSIS

The standard use-case, converting a file and doing something with the meta:

    use Text::Markdown::Frosted qw(markdown);
    use Data::Dump qw(dump);
    
    my $file = 'some/presumed-good/markdown.md';
    my $opts = {
        tab_width     => 2, # passed to Text::Markdown
        keep_tree     => 1, # Frosted option
    };
    my ($meta,$html) = markdown($file, $opts);
    print "$file meta: ", dump($meta);
    print $html;

The slightly more efficient, much more flexible OO usage:

    use Text::Markdown::Frosted;

    my $opts = {
        tab_width     => 2,  # passed to Text::Markdown
        keep_tree     => 1,  # Frosted option, as are:
        expect => {
            title   => 'text',
            author  => 'text',
            tags    => 'list',
            created => 'date',
        },
        expect_strict => 1,
    };
    my $md = Text::Markdown::Frosted->new($opts);
    
    for my $file(@ARGV) {
        my $res = $md->convert($file) or die $md->error;
        process_your_content($res);
    }
    my $html = markdown($text, \%options ); # as in Text::Markdown.

The fallback, compatible with C<Text::Markdown>:

    use Text::Markdown::Frosted qw(markdown);
    my $file = 'some/presumed-good/markdown.md';
    my $opts = {
        tab_width          => 2,  # passed to Text::Markdown
        preserve_html_meta => 1,  # Frosted option
    };
    my $html = markdown($file, $opts);
    print  $html;

Note that the "fallback" is just the standard functional interface used in
scalar mode as opposed to list mode.

=head1 EXPORT

A list of functions that can be exported.  You can delete this section
if you don't export anything, such as for a purely object-oriented module.

=head1 SUBROUTINES/METHODS

=head2 function1

=cut

sub function1 {
}

=head2 function2

=cut

sub function2 {
}

=head1 AUTHOR

Kevin Frost, C<< <biztos at mac.com> >>

=head1 BUGS

Please report any bugs or feature requests to C<bug-text-markdown-frosted at rt.cpan.org>, or through
the web interface at L<http://rt.cpan.org/NoAuth/ReportBug.html?Queue=Text-Markdown-Frosted>.  I will be notified, and then you'll
automatically be notified of progress on your bug as I make changes.




=head1 SUPPORT

You can find documentation for this module with the perldoc command.

    perldoc Text::Markdown::Frosted


You can also look for information at:

=over 4

=item * RT: CPAN's request tracker (report bugs here)

L<http://rt.cpan.org/NoAuth/Bugs.html?Dist=Text-Markdown-Frosted>

=item * AnnoCPAN: Annotated CPAN documentation

L<http://annocpan.org/dist/Text-Markdown-Frosted>

=item * CPAN Ratings

L<http://cpanratings.perl.org/d/Text-Markdown-Frosted>

=item * Search CPAN

L<http://search.cpan.org/dist/Text-Markdown-Frosted/>

=back


=head1 ACKNOWLEDGEMENTS


=head1 LICENSE AND COPYRIGHT

Copyright 2015 Kevin Frost.

This program is free software; you can redistribute it and/or modify it
under the terms of the the Artistic License (2.0). You may obtain a
copy of the full license at:

L<http://www.perlfoundation.org/artistic_license_2_0>

Any use, modification, and distribution of the Standard or Modified
Versions is governed by this Artistic License. By using, modifying or
distributing the Package, you accept this license. Do not use, modify,
or distribute the Package, if you do not accept this license.

If your Modified Version has been derived from a Modified Version made
by someone other than you, you are nevertheless required to ensure that
your Modified Version complies with the requirements of this license.

This license does not grant you the right to use any trademark, service
mark, tradename, or logo of the Copyright Holder.

This license includes the non-exclusive, worldwide, free-of-charge
patent license to make, have made, use, offer to sell, sell, import and
otherwise transfer the Package with respect to any patent claims
licensable by the Copyright Holder that are necessarily infringed by the
Package. If you institute patent litigation (including a cross-claim or
counterclaim) against any party alleging that the Package constitutes
direct or contributory patent infringement, then this Artistic License
to you shall terminate on the date that such litigation is filed.

Disclaimer of Warranty: THE PACKAGE IS PROVIDED BY THE COPYRIGHT HOLDER
AND CONTRIBUTORS "AS IS' AND WITHOUT ANY EXPRESS OR IMPLIED WARRANTIES.
THE IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
PURPOSE, OR NON-INFRINGEMENT ARE DISCLAIMED TO THE EXTENT PERMITTED BY
YOUR LOCAL LAW. UNLESS REQUIRED BY LAW, NO COPYRIGHT HOLDER OR
CONTRIBUTOR WILL BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, OR
CONSEQUENTIAL DAMAGES ARISING IN ANY WAY OUT OF THE USE OF THE PACKAGE,
EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


=cut

1;    # End of Text::Markdown::Frosted
