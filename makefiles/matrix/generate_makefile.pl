#!/usr/bin/env perl

use strict;
use warnings;

my $size = $ARGV[0];
$size = 10 unless defined $size;

my @a;
my @b;
my @c;
my @p;

print "all: c\n\n";
#print "all: c check\n\n";

print "check:\ta b\n";
print "\t multiply check a b\n";

my ($i, $j, $k);
foreach $i (1..$size) {
	foreach $j (1..$size) {
		push @a, "a-$i-$j";
		push @b, "b-$i-$j";
		push @c, "c-$i-$j";
		foreach $k (1..$size) {
			push @p, "p-$i-$k-$k-$j";
		}
	}
}

print "c:\t@c fuse\n";
print "\t./fuse c $size $size @c\n";

print "\n#\n";

my $c;
foreach $c (@c) {
	my @chars = split(/-/, $c);
	my $i = $chars[1];
	my $j = $chars[2];
	my @deps;
	push @deps, "p-$i-$_-$_-$j" foreach (1..$size);
	print "$c:\t@deps sum\n";
	print "\t./sum $c @deps\n";
}

print "\n#\n";

my $p;
foreach $p (@p) {
	my @chars = split(/-/, $p);
	my $i = $chars[1];
	my $k = $chars[2];
	my $j = $chars[4];
	print "$p:\ta-$i-$k b-$k-$j multiply\n";
	print "\t./multiply $p a-$i-$k b-$k-$j\n"
}

print "\n#\n";

foreach $a (@a) {
	my @chars = split(/-/, $a);
	my $i = $chars[1];
	my $j = $chars[2];
	print "$a:\ta split\n";
	print "\t./split $a a $size $size $i $j\n";
}

print "\n#\n";

foreach $b (@b) {
	my @chars = split(/-/, $b);
	my $i = $chars[1];
	my $j = $chars[2];
	print "$b:\tb split\n";
	print "\t./split $b b $size $size $i $j\n";
}

print "clean:\n";
print "\trm -f @a @b @p @c c check\n";
