package strs

// TODO: there is probably a simplification in the logic in this method AllBetweenPattern. We unfortunately can't use a regular expression since go does not support before text matching (?=re) https://github.com/google/re2/wiki/Syntax https://github.com/google/re2/wiki/WhyRE2
func AllBetweenPattern(str, pattern string) []string {
	stringsMatched := []string{}

	matchingPatternIdx := 0
	patternMatched := false

	currentMatch := ""

	for _, char := range str {
		if byte(char) == pattern[matchingPatternIdx] {
			matchingPatternIdx++

			if !patternMatched && matchingPatternIdx == len(pattern) {
				patternMatched = true
				matchingPatternIdx = 0
				continue
			}
		}

		if patternMatched {
			currentMatch += string(char)

			if matchingPatternIdx == len(pattern) {
				patternMatched = false
				matchingPatternIdx = 0

				matchedStr := currentMatch[:len(currentMatch)-len(pattern)]
				stringsMatched = append(stringsMatched, matchedStr)
				currentMatch = ""
			}
		}

	}

	return stringsMatched
}
