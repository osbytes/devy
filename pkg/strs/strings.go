package strs

func AllBetweenPattern(s, pattern string) []string {
	stringsMatched := []string{}

	matchingPatternIdx := 0
	patternMatched := false

	currentMatch := ""

	for _, c := range s {
		if byte(c) == pattern[matchingPatternIdx] {
			matchingPatternIdx++

			if !patternMatched && matchingPatternIdx == len(pattern) {
				patternMatched = true
				matchingPatternIdx = 0
				continue
			}
		}

		if patternMatched {
			currentMatch += string(c)

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
