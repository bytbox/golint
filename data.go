// Copyright 2010 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
		if len(parts) >= 2 {
			key := fmt.Sprintf("deprecated:%s", parts[0])
			value := strings.Join(parts[1:len(parts)], " ")
			statelessLinters[key] = FuncDeprecationLint(parts[0], value)
		}
	}
}

func initMethodDeprecations() {
	lines := strings.Split(methodDeprecations, "\n", -1)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			key := fmt.Sprintf("deprecated:%s.%s", parts[0], parts[1])
			value := strings.Join(parts[2:len(parts)], " ")
			parsingLinters[key] =
				&MethodDeprecationLint{Type: parts[0],
					Method: parts[1], Reason: value}
		}
	}
}

var funcDeprecations = `
new                                     use &T{}
panicln                                 use panic(fmt.Sprintf())
`

var methodDeprecations = `
regexp.Regexp             Execute                   use Find
regexp.Regexp             ExecuteString             use FindString
regexp.Regexp             MatchStrings              use FindStringSubmatch
regexp.Regexp             MatchSlices               use FindSubmatch
regexp.Regexp             AllMatches                use FindAll
regexp.Regexp             AllMatchesString          use FindAllString
`
