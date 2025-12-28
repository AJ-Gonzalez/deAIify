package linter

import (
	"fmt"
	"os/exec"
	"strings"
)

// Tool represents a linting/checking tool
type Tool struct {
	Name       string
	Command    string
	Args       []string
	InstallCmd string
	InstallURL string
	Languages  []string // "js", "py"
}

// Available tools in order of preference
var JSTools = []Tool{
	{
		Name:       "ESLint",
		Command:    "eslint",
		Args:       []string{"--fix"},
		InstallCmd: "npm install -g eslint",
		InstallURL: "https://eslint.org/docs/user-guide/getting-started",
		Languages:  []string{"js"},
	},
	{
		Name:       "Prettier",
		Command:    "prettier",
		Args:       []string{"--write"},
		InstallCmd: "npm install -g prettier",
		InstallURL: "https://prettier.io/docs/en/install.html",
		Languages:  []string{"js"},
	},
}

var GoTools = []Tool{
	{
		Name:       "gofmt",
		Command:    "gofmt",
		Args:       []string{"-w"},
		InstallCmd: "Included with Go installation",
		InstallURL: "https://golang.org/doc/install",
		Languages:  []string{"go"},
	},
	{
		Name:       "staticcheck",
		Command:    "staticcheck",
		Args:       []string{},
		InstallCmd: "go install honnef.co/go/tools/cmd/staticcheck@latest",
		InstallURL: "https://staticcheck.io/",
		Languages:  []string{"go"},
	},
	{
		Name:       "golangci-lint",
		Command:    "golangci-lint",
		Args:       []string{"run"},
		InstallCmd: "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		InstallURL: "https://golangci-lint.run/",
		Languages:  []string{"go"},
	},
}

var PyTools = []Tool{
	{
		Name:       "Ruff",
		Command:    "ruff",
		Args:       []string{"check", "--fix"},
		InstallCmd: "pip install ruff",
		InstallURL: "https://github.com/astral-sh/ruff",
		Languages:  []string{"py"},
	},
	{
		Name:       "Black",
		Command:    "black",
		Args:       []string{},
		InstallCmd: "pip install black",
		InstallURL: "https://black.readthedocs.io/",
		Languages:  []string{"py"},
	},
	{
		Name:       "Flake8",
		Command:    "flake8",
		Args:       []string{},
		InstallCmd: "pip install flake8",
		InstallURL: "https://flake8.pycqa.org/",
		Languages:  []string{"py"},
	},
}

// Syntax checkers (just verify code is valid)
var JSSyntaxCheck = Tool{
	Name:       "Node.js",
	Command:    "node",
	Args:       []string{"--check"},
	InstallCmd: "Download from https://nodejs.org/",
	InstallURL: "https://nodejs.org/",
	Languages:  []string{"js"},
}

var PySyntaxCheck = Tool{
	Name:       "Python",
	Command:    "python",
	Args:       []string{"-m", "py_compile"},
	InstallCmd: "Download from https://python.org/",
	InstallURL: "https://python.org/",
	Languages:  []string{"py"},
}

var GoSyntaxCheck = Tool{
	Name:       "Go",
	Command:    "go",
	Args:       []string{"build", "-o", "/dev/null"},
	InstallCmd: "Download from https://golang.org/",
	InstallURL: "https://golang.org/",
	Languages:  []string{"go"},
}

// Result of running a tool
type Result struct {
	Tool    string
	Success bool
	Output  string
	Error   string
}

// AvailableTools checks which tools are installed
type AvailableTools struct {
	JSLinter     *Tool
	PyLinter     *Tool
	GoLinter     *Tool
	JSSyntax     bool
	PySyntax     bool
	GoSyntax     bool
	MissingTools []Tool
}

// DetectTools checks what linting tools are available
func DetectTools() AvailableTools {
	result := AvailableTools{}

	// Check JS linters
	for i := range JSTools {
		if isAvailable(JSTools[i].Command) {
			result.JSLinter = &JSTools[i]
			break
		} else {
			result.MissingTools = append(result.MissingTools, JSTools[i])
		}
	}

	// Check Python linters
	for i := range PyTools {
		if isAvailable(PyTools[i].Command) {
			result.PyLinter = &PyTools[i]
			break
		} else {
			result.MissingTools = append(result.MissingTools, PyTools[i])
		}
	}

	// Check Go linters
	for i := range GoTools {
		if isAvailable(GoTools[i].Command) {
			result.GoLinter = &GoTools[i]
			break
		} else {
			result.MissingTools = append(result.MissingTools, GoTools[i])
		}
	}

	// Check syntax validators
	result.JSSyntax = isAvailable(JSSyntaxCheck.Command)
	result.PySyntax = isAvailable(PySyntaxCheck.Command)
	result.GoSyntax = isAvailable(GoSyntaxCheck.Command)

	return result
}

// isAvailable checks if a command is in PATH
func isAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// CheckSyntax verifies a file has valid syntax
func CheckSyntax(path string, isPython bool, isGo bool) Result {
	var tool Tool
	if isPython {
		tool = PySyntaxCheck
	} else if isGo {
		tool = GoSyntaxCheck
	} else {
		tool = JSSyntaxCheck
	}

	if !isAvailable(tool.Command) {
		return Result{
			Tool:    tool.Name,
			Success: true, // Skip if not available
			Output:  "syntax check skipped (tool not available)",
		}
	}

	args := append(tool.Args, path)
	cmd := exec.Command(tool.Command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Tool:    tool.Name,
			Success: false,
			Output:  string(output),
			Error:   err.Error(),
		}
	}

	return Result{
		Tool:    tool.Name,
		Success: true,
		Output:  "syntax valid",
	}
}

// RunLinter runs the appropriate linter on a file
func RunLinter(path string, isPython bool, isGo bool, tools AvailableTools) Result {
	var tool *Tool
	if isPython {
		tool = tools.PyLinter
	} else if isGo {
		tool = tools.GoLinter
	} else {
		tool = tools.JSLinter
	}

	if tool == nil {
		return Result{
			Tool:    "none",
			Success: true,
			Output:  "no linter available, skipped",
		}
	}

	args := append(tool.Args, path)
	cmd := exec.Command(tool.Command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Some linters return non-zero for warnings, check output
		outStr := string(output)
		if strings.Contains(strings.ToLower(outStr), "error") {
			return Result{
				Tool:    tool.Name,
				Success: false,
				Output:  outStr,
				Error:   err.Error(),
			}
		}
	}

	return Result{
		Tool:    tool.Name,
		Success: true,
		Output:  string(output),
	}
}

// PrintMissingTools prints suggestions for installing missing tools
func PrintMissingTools(tools AvailableTools, verbose bool) {
	if !verbose || len(tools.MissingTools) == 0 {
		return
	}

	fmt.Println("\nOptional linters not found (install for better results):")

	seen := make(map[string]bool)
	for _, t := range tools.MissingTools {
		if seen[t.Name] {
			continue
		}
		seen[t.Name] = true
		fmt.Printf("  %s: %s\n", t.Name, t.InstallCmd)
	}
	fmt.Println()
}
