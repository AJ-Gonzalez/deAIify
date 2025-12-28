# deaiify

A CLI tool that transforms AI-generated code to appear more human-written. Runs locally, no LLM calls, pure pattern matching.

## Installation

```bash
go build -o deaiify ./cmd/deaiify
```

## Usage

```bash
deaiify ./src                    # Process all supported files in directory
deaiify ./src --dry-run          # Show what would change without modifying
deaiify ./main.py                # Process single file
deaiify ./src --verbose          # Show detailed transformation log
deaiify ./src --lint             # Run linters after transformation
deaiify --scan-commits           # Scan git commits for AI patterns
```

## Supported Languages

- TypeScript/JavaScript (.ts, .tsx, .js, .jsx)
- Python (.py)
- Go (.go)

## What It Does

### Comment Humanization

Detects AI patterns like:
- Comments starting with "This function...", "Here's how...", "Let's...", "We need to..."
- Over-explanatory comments (>100 chars for simple operations)
- Overly formal language ("in order to", "this ensures that")

Replaces with human-style comments:
- `// TODO`, `// FIXME`, `// XXX`
- `// hack`, `// fix later`, `// works somehow`
- `// don't touch`, `// legacy`, `// ugh`
- Or removes them entirely (humans under-comment)

### Typo Injection

~5% of comments get a realistic typo:
- "the" -> "teh"
- "that" -> "taht"
- "receive" -> "recieve"
- "definitely" -> "definately"
- Random character swaps

### Formatting Inconsistencies

- Occasionally removes trailing commas
- Adds random extra blank lines
- Minor spacing variations

### Linting Integration

Use `--lint` to run linters after transformation. Supports:

**JavaScript/TypeScript:**
- ESLint (`npm install -g eslint`)
- Prettier (`npm install -g prettier`)

**Python:**
- Ruff (`pip install ruff`)
- Black (`pip install black`)
- Flake8 (`pip install flake8`)

**Go:**
- gofmt (included with Go)
- staticcheck (`go install honnef.co/go/tools/cmd/staticcheck@latest`)
- golangci-lint (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)

Linters are auto-detected. If none are found, syntax checking still runs via Node.js/Python to ensure valid code.

### Git Commit Scanning

Scans recent commits and warns about:
- Emoji in commit messages
- "Generated with" footers
- "Co-Authored-By: Claude" or similar AI footers
- Overly formal commit messages

## Examples

Before:
```javascript
// This function handles the user authentication process
// Here's how we validate the user credentials
function authenticateUser(username, password) {
  // Let's first check if the username is valid
  if (!username) {
    // This ensures that we have a valid username before proceeding
    return false;
  }
  // ...
}
```

After:
```javascript
// hack
function authenticateUser(username, password) {
  // XXX
  if (!username) {
    return false;
  }
  // ...
}
```

## Non-Goals

- No LLM calls
- No network access
- No semantic understanding beyond pattern matching
- No breaking valid syntax

## Irony Level

Maximum. This tool hides that AI built your code. The tool itself was built by AI.
  
