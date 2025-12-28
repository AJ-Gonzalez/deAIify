package detector

import (
	"regexp"
	"strings"
	"unicode"

	"deaiify/internal/parser"
)

// AI comment patterns - phrases that AI tends to use
var AIPatternPrefixes = []string{
	"this function",
	"this method",
	"this class",
	"this module",
	"this component",
	"this hook",
	"this utility",
	"this helper",
	"here's",
	"here we",
	"here is",
	"let's",
	"let us",
	"i'll",
	"i will",
	"we'll",
	"we will",
	"we need to",
	"we can",
	"the following",
	"below is",
	"above is",
	"as you can see",
	"note that",
	"notice that",
	"importantly",
	"essentially",
	"basically",
}

// Emoji detection regex
var emojiRe = regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F1E0}-\x{1F1FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)

// DetectionResult describes why a comment was flagged
type DetectionResult struct {
	IsAILike    bool
	Reasons     []string
	Score       float64 // 0.0 to 1.0, how "AI-like"
}

// DetectAIComment checks if a comment looks AI-generated
func DetectAIComment(comment parser.Comment) DetectionResult {
	result := DetectionResult{
		Reasons: []string{},
	}
	text := strings.ToLower(strings.TrimSpace(comment.Text))

	// Check for AI prefix patterns
	for _, prefix := range AIPatternPrefixes {
		if strings.HasPrefix(text, prefix) {
			result.Reasons = append(result.Reasons, "starts with AI pattern: "+prefix)
			result.Score += 0.4
			break
		}
	}

	// Check for emoji
	if emojiRe.MatchString(comment.Text) {
		result.Reasons = append(result.Reasons, "contains emoji")
		result.Score += 0.3
	}

	// Check for over-explanation (>100 chars for simple statements)
	if len(comment.Text) > 100 && !comment.IsBlock {
		result.Reasons = append(result.Reasons, "overly verbose single-line comment")
		result.Score += 0.2
	}

	// Check for overly formal language
	if containsFormalLanguage(text) {
		result.Reasons = append(result.Reasons, "uses overly formal language")
		result.Score += 0.15
	}

	// Check for excessive capitalization in explanations
	if hasExcessiveCapitalization(comment.Text) {
		result.Reasons = append(result.Reasons, "excessive capitalization")
		result.Score += 0.1
	}

	// Normalize score
	if result.Score > 1.0 {
		result.Score = 1.0
	}

	result.IsAILike = result.Score >= 0.3

	return result
}

// containsFormalLanguage checks for overly formal phrasing
func containsFormalLanguage(text string) bool {
	formalPatterns := []string{
		"in order to",
		"it is important to",
		"it should be noted",
		"this ensures that",
		"this allows us to",
		"this enables",
		"this provides",
		"this handles",
		"responsible for",
		"utilized",
		"implement",
		"functionality",
	}

	for _, pattern := range formalPatterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}

// hasExcessiveCapitalization checks for unusual capitalization patterns
func hasExcessiveCapitalization(text string) bool {
	if len(text) < 20 {
		return false
	}
	upperCount := 0
	letterCount := 0
	for _, r := range text {
		if unicode.IsLetter(r) {
			letterCount++
			if unicode.IsUpper(r) {
				upperCount++
			}
		}
	}
	if letterCount == 0 {
		return false
	}
	// If more than 30% is uppercase (excluding first letter patterns), flag it
	ratio := float64(upperCount) / float64(letterCount)
	return ratio > 0.3
}
