package transformer

import (
	"math/rand"
	"sort"
	"strings"

	"deaiify/internal/parser"
)

// TypoMap contains common typos humans make
var TypoMap = map[string]string{
	"the":         "teh",
	"that":        "taht",
	"receive":     "recieve",
	"separate":    "seperate",
	"definitely":  "definately",
	"occurred":    "occured",
	"necessary":   "neccessary",
	"until":       "untill",
	"beginning":   "begining",
	"believe":     "beleive",
	"their":       "thier",
	"weird":       "wierd",
	"which":       "wich",
	"because":     "becuase",
	"successful":  "successfull",
	"apparently":  "apparantly",
	"immediately": "immediatly",
	"occurrence":  "occurence",
	"recommend":   "reccomend",
	"reference":   "refrence",
	"calendar":    "calender",
	"environment": "enviroment",
	"available":   "availble",
	"development": "developement",
	"function":    "fucntion",
	"return":      "retrun",
	"variable":    "vairable",
	"parameter":   "paramter",
	"argument":    "arguement",
}

// TypoProbability is the chance a comment gets a typo (5%)
const TypoProbability = 0.05

// InjectTypos randomly adds typos to comments (not code!)
func InjectTypos(content string, comments []parser.Comment, p parser.Parser) (string, []TransformResult) {
	var results []TransformResult

	// Sort by position in descending order
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Start > comments[j].Start
	})

	// Process in reverse position order
	for _, comment := range comments {

		// Skip with 95% probability
		if rand.Float64() > TypoProbability {
			continue
		}

		// Try to inject a typo
		newText, typoMade := injectTypoInText(comment.Text)
		if !typoMade {
			continue
		}

		content = p.ReplaceComment(content, comment, newText)
		results = append(results, TransformResult{
			Original:    comment.Original,
			Replacement: newText,
			Action:      "typo_injected",
			LineNumber:  comment.LineNumber,
		})
	}

	return content, results
}

// injectTypoInText attempts to add a typo to the text
func injectTypoInText(text string) (string, bool) {
	lower := strings.ToLower(text)

	// Try each word in the typo map
	for correct, typo := range TypoMap {
		if idx := strings.Index(lower, correct); idx != -1 {
			// Preserve original case of first letter
			typoCased := typo
			if idx < len(text) && text[idx] >= 'A' && text[idx] <= 'Z' {
				typoCased = strings.ToUpper(string(typo[0])) + typo[1:]
			}

			// Replace the word
			result := text[:idx] + typoCased + text[idx+len(correct):]
			return result, true
		}
	}

	// If no dictionary typo found, try random character swap (20% chance)
	if rand.Float64() < 0.2 && len(text) > 10 {
		return randomCharSwap(text)
	}

	return text, false
}

// randomCharSwap swaps two adjacent characters
func randomCharSwap(text string) (string, bool) {
	runes := []rune(text)

	// Find a good position (in the middle of a word)
	attempts := 0
	for attempts < 5 {
		pos := rand.Intn(len(runes) - 1)
		if isSwappable(runes[pos]) && isSwappable(runes[pos+1]) {
			runes[pos], runes[pos+1] = runes[pos+1], runes[pos]
			return string(runes), true
		}
		attempts++
	}

	return text, false
}

// isSwappable returns true if the character is a lowercase letter
func isSwappable(r rune) bool {
	return r >= 'a' && r <= 'z'
}
