package transformer

import (
	"math/rand"
	"regexp"
	"strings"
)

// StructureTransformResult describes structural changes made
type StructureTransformResult struct {
	Type        string // "trailing_comma", "blank_line", "quote_style", etc.
	Description string
	LineNumber  int
}

// TransformStructure applies formatting inconsistencies
func TransformStructure(content string, isJS bool) (string, []StructureTransformResult) {
	var results []StructureTransformResult

	if isJS {
		content, results = transformJSStructure(content, results)
	} else {
		content, results = transformPyStructure(content, results)
	}

	return content, results
}

// transformJSStructure applies JS/TS specific transformations
func transformJSStructure(content string, results []StructureTransformResult) (string, []StructureTransformResult) {
	lines := strings.Split(content, "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Randomly remove trailing comma (10% chance)
		if rand.Float64() < 0.1 {
			if newLine, changed := removeTrailingComma(line); changed {
				lines[i] = newLine
				results = append(results, StructureTransformResult{
					Type:        "trailing_comma",
					Description: "removed trailing comma",
					LineNumber:  i + 1,
				})
			}
		}

		// Randomly add extra blank line after closing brace (5% chance)
		if rand.Float64() < 0.05 && strings.TrimSpace(line) == "}" {
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) != "" {
				lines = insertLine(lines, i+1, "")
				results = append(results, StructureTransformResult{
					Type:        "blank_line",
					Description: "added extra blank line",
					LineNumber:  i + 1,
				})
				i++ // Skip the inserted line
			}
		}

		// Operator spacing transformation disabled - too risky for code correctness
		// Keeping code safe is more important than humanization
	}

	return strings.Join(lines, "\n"), results
}

// transformPyStructure applies Python specific transformations
func transformPyStructure(content string, results []StructureTransformResult) (string, []StructureTransformResult) {
	lines := strings.Split(content, "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Randomly add extra blank line after function def (5% chance)
		if rand.Float64() < 0.05 && strings.HasPrefix(strings.TrimSpace(line), "def ") {
			// Check if next line isn't already blank
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) != "" {
				lines = insertLine(lines, i+1, "")
				results = append(results, StructureTransformResult{
					Type:        "blank_line",
					Description: "added blank line after function",
					LineNumber:  i + 1,
				})
				i++
			}
		}

		// Randomly remove trailing comma in lists/dicts (10% chance)
		if rand.Float64() < 0.1 {
			if newLine, changed := removeTrailingComma(line); changed {
				lines[i] = newLine
				results = append(results, StructureTransformResult{
					Type:        "trailing_comma",
					Description: "removed trailing comma",
					LineNumber:  i + 1,
				})
			}
		}
	}

	return strings.Join(lines, "\n"), results
}

// removeTrailingComma removes trailing comma before ) or ]
func removeTrailingComma(line string) (string, bool) {
	trimmed := strings.TrimRight(line, " \t")
	patterns := []string{",)", ",]", ",}"}

	for _, p := range patterns {
		if strings.HasSuffix(trimmed, p) {
			return strings.TrimSuffix(trimmed, p) + p[1:], true
		}
	}
	return line, false
}

// tweakOperatorSpacing randomly adjusts spacing around operators
func tweakOperatorSpacing(line string) (string, bool) {
	// Only target simple assignments not in strings
	// This is a simplified heuristic

	// Pattern: adjust spacing around = (not == or !=)
	re := regexp.MustCompile(`(\w)\s*=\s*(\w)`)
	if re.MatchString(line) {
		// Randomly pick a style
		styles := []string{"$1=$2", "$1 = $2", "$1= $2", "$1 =$2"}
		style := styles[rand.Intn(len(styles))]

		// Only change once per line
		newLine := re.ReplaceAllString(line, style)
		if newLine != line {
			return newLine, true
		}
	}

	return line, false
}

// insertLine inserts a new line at the given position
func insertLine(lines []string, pos int, newLine string) []string {
	result := make([]string, len(lines)+1)
	copy(result[:pos], lines[:pos])
	result[pos] = newLine
	copy(result[pos+1:], lines[pos:])
	return result
}
