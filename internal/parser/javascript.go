package parser

import (
	"regexp"
	"strings"
)

// JavaScriptParser parses JS/TS files for comments
type JavaScriptParser struct{}

// NewJavaScriptParser creates a new JavaScript parser
func NewJavaScriptParser() *JavaScriptParser {
	return &JavaScriptParser{}
}

var (
	// Match single-line comments: // ...
	jsLineCommentRe = regexp.MustCompile(`(?m)//[^\n]*`)
	// Match block comments: /* ... */
	jsBlockCommentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// Parse extracts all comments from JavaScript/TypeScript content
func (p *JavaScriptParser) Parse(content string) ParseResult {
	var comments []Comment

	// Find line comments
	for _, match := range jsLineCommentRe.FindAllStringIndex(content, -1) {
		start, end := match[0], match[1]
		original := content[start:end]
		text := strings.TrimPrefix(original, "//")
		text = strings.TrimSpace(text)
		lineNum := countLines(content[:start]) + 1

		comments = append(comments, Comment{
			Text:       text,
			Start:      start,
			End:        end,
			LineNumber: lineNum,
			IsBlock:    false,
			Original:   original,
		})
	}

	// Find block comments
	for _, match := range jsBlockCommentRe.FindAllStringIndex(content, -1) {
		start, end := match[0], match[1]
		original := content[start:end]
		text := strings.TrimPrefix(original, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)
		lineNum := countLines(content[:start]) + 1

		comments = append(comments, Comment{
			Text:       text,
			Start:      start,
			End:        end,
			LineNumber: lineNum,
			IsBlock:    true,
			Original:   original,
		})
	}

	return ParseResult{
		Content:  content,
		Comments: comments,
	}
}

// ReplaceComment replaces a comment in the content with new text
func (p *JavaScriptParser) ReplaceComment(content string, comment Comment, newText string) string {
	var replacement string
	if newText == "" {
		// Remove the comment entirely
		startIdx := comment.Start
		endIdx := comment.End

		// If comment is at start of line (only whitespace before it), remove the whole line
		lineStart := startIdx
		for lineStart > 0 && content[lineStart-1] != '\n' {
			lineStart--
		}
		prefix := content[lineStart:startIdx]
		if strings.TrimSpace(prefix) == "" {
			// Only whitespace before comment, remove from line start
			startIdx = lineStart
		}

		// Remove trailing newline if present
		if endIdx < len(content) && content[endIdx] == '\n' {
			endIdx++
		}

		return content[:startIdx] + content[endIdx:]
	}

	if comment.IsBlock {
		replacement = "/* " + newText + " */"
	} else {
		replacement = "// " + newText
	}
	return content[:comment.Start] + replacement + content[comment.End:]
}

func countLines(s string) int {
	return strings.Count(s, "\n")
}
