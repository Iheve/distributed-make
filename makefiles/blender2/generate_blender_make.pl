#!/usr/bin/env perl
use strict;
use warnings;
use List::Util qw(sum);

my @frames;
my @files;

for my $arg (@ARGV) {
	if ($arg =~ /.blend$/) {
		push @files, $arg;
	} else {
		push @frames, $arg;
	}
}

die "incorrect args" unless $#frames == $#files;
#read this file from last to first line to make sense of it

#generate video
my $c = (sum @frames) - 9 * (@files -1);
my @images = map {"f_$_.jpg"} (1..$c);
my $images = join(' ', @images);
print "out.avi: $images\n";
print "\t".'ffmpeg -i "f_%d.jpg" out.avi'."\n\n";

#convert and rename all images (not in transitions)
my $count = 0;
for my $idx (0..$#files) {
	for my $num (1..$frames[$idx]) {
		#skip firt 9 images unless first file
		if ($idx != 0) {
			if ($num <= 9) {
				next;
			}
		}
		$count++;
		#skip last 9 images unless last file
		if ($idx != $#files) {
			if ($num >= $frames[$idx] - 8) {
				next;
			}
		}
		my $file = $files[$idx];
		$file =~ /([^\/]+)\.blend$/;
		my $in = "$1_$num.tga";
		my $out = "f_$count.jpg";
		print "$out: $in\n";
		print "\tconvert $in -resize 640x480 $out\n\n";
	}
}



#transition images
$count = 0;
for my $idx (1..$#files) {
	$count += $frames[$idx-1] - 9;
	for my $num (1..9) {
		my $n = $count + $num;
		my $out = "f_$n.jpg";
		my $blend1 = $files[$idx-1];
		$blend1 =~ /([^\/]+).blend$/;
		my $in1 = $1."_".($frames[$idx-1] - $num).".tga";
		my $blend2 = $files[$idx];
		$blend2 =~ /([^\/]+).blend$/;
		my $in2 = $1."_".($num).".tga";
		print "$out: $in1 $in2\n";
		print "\tcomposite -resize 640x480 $in1 $in2 -blend ".((10-$num)*10)."x".($num*10)." $out\n\n";
	}
}

#images computations
for my $idx (0..$#files) {
	for my $num (1..$frames[$idx]) {
		my $in = $files[$idx];
		$in =~ /([^\/]+).blend$/;
		my $root = $1;
		my $output = "$1_$num.tga";
		print "$output : $in\n";
		print "\tblender -b $in -o //${root}_#.tga -F TGA -f $num\n\n";
	}
}


