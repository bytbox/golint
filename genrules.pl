#!/usr/bin/perl

open RULES, ">rules.go";

print RULES <<END;
package main

var LineLinters = [...]LineLinter{
END

open FIN, "rules/line-regex";
while (<FIN>) {
	$name = $_;
	chomp $name;
	$regex = <FIN> or break;
	chomp $regex;
	$regex =~ s/^\t//;
	$desc = <FIN> or break;
	chomp $desc;
	$desc =~ s/^\t//;

	@nameparts = split /:/, $name;
	$category = $nameparts[0];
	$name2 = $nameparts[1];
	$name2 =~ s/-/_/g;

print RULES <<END;
RegexLinter{LinterName{"$category", "$name2", "$desc"}, `$regex`},
END

	<FIN>;
}

print RULES <<END;
}
END

close FIN;
close RULES;

