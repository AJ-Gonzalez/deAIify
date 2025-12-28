package parser

// Comment represents a comment found in source code
type Comment struct {
	Text       string // The comment text (without delimiters)
	Start      int    // Start position in the file
	End        int    // End position in the file
	LineNumber int    // Line number (1-indexed)
	IsBlock    bool   // True if block comment (/* */ or """)
	Original   string // Original text including delimiters
}

// ParseResult holds the result of parsing a file
type ParseResult struct {
	Content  string
	Comments []Comment
}

// Parser interface for language-specific parsers
type Parser interface {
	Parse(content string) ParseResult
	ReplaceComment(content string, comment Comment, newText string) string
}
