package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	clearSpecialCharEnding    = regexp.MustCompile(`[!?;:,.'"*+<>=#%&/(){}~|\[\]\\]\s`)
	clearSpecialCharBeginning = regexp.MustCompile(`\s[!?;:,.'"*+<>=#%&/(){}~|\[\]\\]`)
	clearSingleDash           = regexp.MustCompile(`\s-\s`)
)

func Top10(t string) []string {
	t = clearSpecialCharBeginning.ReplaceAllLiteralString(strings.ToLower(t), " ")
	t = clearSpecialCharEnding.ReplaceAllLiteralString(t, " ")
	t = clearSingleDash.ReplaceAllLiteralString(t, "")

	words := strings.Fields(t)
	freq := make(map[string]int, len(words))

	for w := range words {
		freq[words[w]]++
	}

	freqInverse := make(map[int][]string)

	for k, v := range freq {
		freqInverse[v] = append(freqInverse[v], k)
	}

	res := make([]string, 0, 10)
	for f := len(t) - 1; f > 0; f-- {
		sort.Strings(freqInverse[f])
		for w := range freqInverse[f] {
			if len(res) == 10 {
				break
			}
			res = append(res, freqInverse[f][w])
		}
	}
	return res
}
