package main

import (
	"fmt"
	"strings"
)

func initLints() {
	initFuncDeprecations()
	initMethodDeprecations()
}

func initFuncDeprecations() {
	lines := strings.Split(funcDeprecations, "\n", -1)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts)>=2 {
			key := fmt.Sprintf("deprecated:%s", parts[0])
			value := strings.Join(parts[1:len(parts)], " ")
			statelessLinters[key] = FuncDeprecationLint(parts[0], value)
		}
	}
}

func initMethodDeprecations() {
	lines := strings.Split(funcDeprecations, "\n", -1)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts)==2 {
			// TODO methodDeprecations
		}
	}
}

var funcDeprecations=
`
new                                     use &T{}
`

var methodDeprecations=
`
regexp.Regexp             Execute                   use Find
regexp.Regexp             ExecuteString             use FindString
regexp.Regexp             MatchStrings              use FindStringSubmatch
regexp.Regexp             MatchSlices               use FindSubmatch
regexp.Regexp             AllMatches                use FindAll
regexp.Regexp             AllMatchesString          use FindAllString
`
