# deaiify - CLI Tool Specification

## Purpose
A CLI tool that transforms AI-generated code to appear more human-written. Runs locally, no LLM calls, pure pattern matching and transformation.

## Usage
```bash
deaiify ./src                    # Process all supported files in directory
deaiify ./src --dry-run          # Show what would change without modifying
deaiify ./main.py                # Process single file
deaiify ./src --verbose          # Show detailed transformation log
```

## Supported Languages
- TypeScript/JavaScript (.ts, .tsx, .js, .jsx)
- Python (.py)

## Transformations

### 1. Comment Humanization
**Detect & Remove AI patterns:**
- Comments starting with "This", "Here's", "Let's", "I'll", "We"
- Over-explanatory comments (>100 chars explaining simple operations)
- Comments that restate the code literally
- Any emoji in comments

**Replace with human patterns:**
- Terse comments: "// handle edge case", "# TODO", "// fix later"
- Occasionally remove comments entirely (humans under-comment)
- Add sporadic: "// don't touch this", "// works somehow", "// legacy"

### 2. Typo Injection (Comments Only)
- Random character swaps: "teh", "taht", "recieve"
- Occasional missing letters in longer words
- Probability: ~5% of comments get a typo
- Never in code, only in comments/strings meant for humans

### 3. Formatting Inconsistencies
- Occasionally skip trailing comma in last array/object item
- Mix quote styles slightly in JS/TS (but keep it parseable)
- Random extra blank lines in some places
- Inconsistent spacing around operators (where linter won't catch)

### 4. Variable/Function Name Variations
- If everything is perfectly camelCase, occasionally use slightly awkward names
- Add a `temp`, `data2`, `handleClick2` somewhere
- Keep a "mistake" variable that's declared but unused

### 5. Structure De-optimization
- Occasionally inline a small function that "should" be extracted
- Leave a slightly repetitive code block instead of DRY-ing it
- Add a `// TODO: refactor this` above it for authenticity

### 6. Commit Message Cleaning (if .git detected)
- Flag emoji in recent commit messages
- Suggest rewrites without emoji
- Remove "Generated with" or "Co-Authored-By: Claude" footers

## Architecture

```
cmd/
  deaiify/
    main.go           # CLI entry, flag parsing
internal/
  detector/
    patterns.go       # Regex patterns for AI detection
    scorer.go         # Score how "AI-like" a file is
  transformer/
    comments.go       # Comment transformation logic
    typos.go          # Typo injection
    structure.go      # Code structure changes
  parser/
    javascript.go     # JS/TS parsing (use tree-sitter or regex)
    python.go         # Python parsing
  walker/
    files.go          # Directory traversal, file filtering
```

## Output
- Modified files in place (or --dry-run for preview)
- Summary: "Processed 23 files, made 47 transformations"
- Verbose mode shows each transformation

## Non-Goals
- No LLM calls
- No network access
- No semantic understanding beyond pattern matching
- No breaking valid syntax

## Irony Level
Maximum. This tool is built to hide that AI built your code. The tool itself may be built by AI.
