#!/usr/bin/perl

sub getNameparts {
	$name = shift;
	chomp $name;
	@nameparts = split /:/, $name;
	$category = $nameparts[0];
	$name = $nameparts[1];
	$name2 = $name;
	$name2 =~ s/-/_/g;
	return ($category, $name, $name2)
}

open RULES, ">rules.go";

print RULES <<END;
package main

var LineLinters = [...]LineLinter{
END

open FIN, "rules/line-regex";
while (<FIN>) {
	$name = $_;
	$regex = <FIN> or break;
	chomp $regex;
	$regex =~ s/^\t//;
	$desc = <FIN> or break;
	chomp $desc;
	$desc =~ s/^\t//;
	($category, $name, $name2) = getNameparts $name;
	print RULES <<END;
RegexLinter{LinterName{"$category", "$name", "$desc"}, `$regex`},
END

	<FIN>;
}
close FIN;

opendir DIR, "rules/line-simple" or die "Could not read rules/line-simple: $!";
while ($fname = readdir(DIR)) {
	next if $fname =~ /^\./;
	open FIN, "rules/line-simple/$fname" or die "Could not open $fname: $!";
	$name = <FIN>;
	($category, $name, $name2) = getNameparts $name;
	$desc = <FIN>;
	chomp $desc;
	# read the rest of the file as code
	$code = "";
	while ($line = <FIN>) {
		$code .= $line;
	}
	$code =~ s/\n+$//msg;
	print RULES <<END;
SimpleLineLinter{LinterName{"$category", "$name", "$desc"},
$code},
END
	close FIN;
}
closedir DIR;

print RULES <<END;
}
END

close RULES;

