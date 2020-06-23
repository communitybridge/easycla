// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"regexp"
)

// ValidCompanyName is a routine to indicate if the regex is a valid company name
func ValidCompanyName(name string) bool {
	// The {2,255} max limit isn't working in the regex, so, use this approach
	if len(name) < 2 || len(name) > 255 {
		return false
	}

	// \w represents any character from the class [A-Za-z0-9_], mnemonic: 'word'.
	// \s represents any of the following whitespaces: TAB, SPACE, CR, LF. Or more precisely [\t\n\f\r ].
	//r := regexp.MustCompile(`^([\w\u00C0-\u00FF][\w\s\u00C0-\u00FF()\[\]+\-/%!@#$]*){2,255}$`)
	// \p{L} represents the latin characters
	r := regexp.MustCompile(`^([\w\p{L}][\w\s\p{L}()\[\]+\-/%!@#$]*){2,255}$`)
	return r.MatchString(name)
}
