#!/usr/bin/env perl

use strict;
use warnings;

my ($w, $h) = @ARGV;
die "args missing" unless defined $h;

print "$w $h\n";
foreach(1..$w) {
	my @line;
	foreach(1..$h) {
		push @line, int rand(10)+1;
	}
	print "@line\n";
}

