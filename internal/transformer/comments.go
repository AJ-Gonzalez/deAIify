package transformer

import (
	"math/rand"
	"sort"

	"deaiify/internal/detector"
	"deaiify/internal/parser"
)

// HumanComments are terse, realistic comments humans actually write
var HumanComments = []string{
	"TODO",
	"FIXME",
	"XXX",
	"hack",
	"fix later",
	"works somehow",
	"don't touch",
	"legacy",
	"ugh",
	"temp",
	"cleanup needed",
	"wtf",
	"why?",
	"magic number",
	"sorry",
	"good enough",
	"needs refactor",
	"not ideal",
	"idk",
}

// TransformResult describes what transformation was applied
type TransformResult struct {
	Original    string
	Replacement string
	Action      string // "removed", "replaced", "typo_injected", "kept"
	LineNumber  int
}

// TransformComments processes AI-detected comments and humanizes them
func TransformComments(content string, p parser.Parser, scores []detector.CommentScore, isPython bool) (string, []TransformResult) {
	var results []TransformResult

	// Sort by position (Start) in descending order so we process from end to start
	// This ensures position validity as we modify the content
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Comment.Start > scores[j].Comment.Start
	})

	// Process comments in reverse position order
	for _, cs := range scores {
		if !cs.Result.IsAILike {
			continue
		}

		action := decideAction(cs)
		var newContent string
		var result TransformResult

		result.Original = cs.Comment.Original
		result.LineNumber = cs.Comment.LineNumber

		switch action {
		case "remove":
			newContent = p.ReplaceComment(content, cs.Comment, "")
			result.Action = "removed"
			result.Replacement = ""

		case "replace":
			replacement := pickHumanComment(isPython)
			newContent = p.ReplaceComment(content, cs.Comment, replacement)
			result.Action = "replaced"
			result.Replacement = replacement

		default:
			result.Action = "kept"
			newContent = content
		}

		content = newContent
		results = append(results, result)
	}

	return content, results
}

// decideAction randomly decides what to do with an AI comment
func decideAction(cs detector.CommentScore) string {
	r := rand.Float64()

	// Higher score = more aggressive transformation
	if cs.Result.Score > 0.7 {
		if r < 0.6 {
			return "remove" // humans under-comment
		}
		return "replace"
	}

	if cs.Result.Score > 0.4 {
		if r < 0.4 {
			return "remove"
		}
		if r < 0.8 {
			return "replace"
		}
		return "keep"
	}

	// Lower scores - be more conservative
	if r < 0.3 {
		return "remove"
	}
	if r < 0.6 {
		return "replace"
	}
	return "keep"
}

// pickHumanComment selects a random human-style comment
func pickHumanComment(isPython bool) string {
	comment := HumanComments[rand.Intn(len(HumanComments))]

	// Occasionally make it lowercase or add punctuation
	if rand.Float64() < 0.3 {
		comment = comment + "..."
	}

	return comment
}

// shortenComment reduces a verbose comment to just a few words
func shortenComment(text string) string {
	words := splitWords(text)
	if len(words) <= 3 {
		return text
	}

	// Take first 2-4 words
	n := 2 + rand.Intn(3)
	if n > len(words) {
		n = len(words)
	}

	result := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			result += " "
		}
		result += words[i]
	}
	return result
}

// splitWords splits text into words
func splitWords(text string) []string {
	var words []string
	var current string

	for _, r := range text {
		if r == ' ' || r == '\t' || r == '\n' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
