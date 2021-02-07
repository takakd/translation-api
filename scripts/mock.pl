#!/usr/bin/perl

# Usage
sub usage() {
    my $usage= <<"EOS";
$0 creates go mock file in the same directory where the go file is located.

Usage:
  $0 <go file path>
EOS
    print($usage);
}

# Check arg count.
my $argc = @ARGV;
if ($argc < 1) {
    usage;
    exit 1;
}

# Get go file.
my $file = $ARGV[0];

# Get package name in the input go file.
my $pkg = "";

open(FH, '<', $file) or die $!;
while(<FH>){
    # Parse package definition line.
    if ($_ =~ /^package ([a-z_0-9]+)$/) {
        $pkg = $1;
        last;
    }
}
close(FH);

# Add "_mock" to the output mock file.
my $mock_file = $file;
$mock_file =~ s/.go$/_mock.go/;

# mockgen command.
my $cmd = "mockgen -source=${file} -destination=${mock_file} -package=${pkg}";

# Run command.
system($cmd);

