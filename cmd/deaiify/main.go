package main

import (
	"flag"
	"fmt"
	"os"

	"deaiify/internal/detector"
	"deaiify/internal/git"
	"deaiify/internal/parser"
	"deaiify/internal/transformer"
	"deaiify/internal/walker"
)

var (
	dryRun      = flag.Bool("dry-run", false, "Show what would change without modifying files")
	verbose     = flag.Bool("verbose", false, "Show detailed transformation log")
	scanCommits = flag.Bool("scan-commits", false, "Scan git commits for AI patterns")
	commitCount = flag.Int("commits", 20, "Number of commits to scan (default 20)")
)

func main() {
	flag.Parse()

	// Handle --scan-commits mode
	if *scanCommits {
		path := "."
		if flag.NArg() > 0 {
			path = flag.Arg(0)
		}
		runCommitScan(path)
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("Usage: deaiify <path> [options]")
		fmt.Println("\nTransforms AI-generated code to appear more human-written.")
		fmt.Println("\nOptions:")
		fmt.Println("  --dry-run       Show what would change without modifying files")
		fmt.Println("  --verbose       Show detailed transformation log")
		fmt.Println("  --scan-commits  Scan git commits for AI patterns")
		fmt.Println("  --commits N     Number of commits to scan (default 20)")
		os.Exit(1)
	}

	path := flag.Arg(0)

	// Validate path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: path does not exist: %s\n", path)
		os.Exit(1)
	}

	// Walk the path to find files
	files, err := walker.Walk(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking path: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No supported files found.")
		os.Exit(0)
	}

	if *verbose {
		fmt.Printf("Found %d files to process\n\n", len(files))
	}

	totalFiles := 0
	totalTransformations := 0

	for _, file := range files {
		transformations, err := processFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file.Path, err)
			continue
		}

		if transformations > 0 {
			totalFiles++
			totalTransformations += transformations
		}
	}

	// Summary
	if *dryRun {
		fmt.Printf("\n[DRY RUN] Would process %d files, make %d transformations\n", totalFiles, totalTransformations)
	} else {
		fmt.Printf("\nProcessed %d files, made %d transformations\n", totalFiles, totalTransformations)
	}
}

func processFile(file walker.FileInfo) (int, error) {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return 0, err
	}

	originalContent := string(content)
	currentContent := originalContent

	// Select parser
	var p parser.Parser
	isPython := file.IsPython()
	if isPython {
		p = parser.NewPythonParser()
	} else {
		p = parser.NewJavaScriptParser()
	}

	// Parse comments
	result := p.Parse(currentContent)

	// Score comments
	score := detector.ScoreFile(file.Path, result.Comments)

	if *verbose && len(score.GetAIComments()) > 0 {
		fmt.Printf("File: %s (AI score: %.2f)\n", file.Path, score.Score)
	}

	var allResults []transformer.TransformResult

	// Transform AI-detected comments
	aiComments := score.GetAIComments()
	if len(aiComments) > 0 {
		newContent, results := transformer.TransformComments(currentContent, p, aiComments, isPython)
		currentContent = newContent
		allResults = append(allResults, results...)
	}

	// Inject typos in remaining comments (re-parse after transforms)
	result = p.Parse(currentContent)
	newContent, typoResults := transformer.InjectTypos(currentContent, result.Comments, p)
	currentContent = newContent
	allResults = append(allResults, typoResults...)

	// Apply structural changes
	newContent, structResults := transformer.TransformStructure(currentContent, !isPython)
	currentContent = newContent

	transformCount := len(allResults) + len(structResults)

	// Print verbose output
	if *verbose {
		for _, r := range allResults {
			fmt.Printf("  Line %d: %s\n", r.LineNumber, r.Action)
			if r.Action != "removed" && r.Replacement != "" {
				fmt.Printf("    -> %s\n", r.Replacement)
			}
		}
		for _, r := range structResults {
			fmt.Printf("  Line %d: %s\n", r.LineNumber, r.Description)
		}
		if transformCount > 0 {
			fmt.Println()
		}
	}

	// Write if not dry run and there are changes
	if !*dryRun && currentContent != originalContent {
		err = os.WriteFile(file.Path, []byte(currentContent), 0644)
		if err != nil {
			return 0, err
		}
	}

	return transformCount, nil
}

func runCommitScan(path string) {
	if !git.IsGitRepo(path) {
		fmt.Fprintf(os.Stderr, "Error: %s is not a git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Scanning last %d commits for AI patterns...\n\n", *commitCount)

	warnings, err := git.ScanCommits(path, *commitCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning commits: %v\n", err)
		os.Exit(1)
	}

	if len(warnings) == 0 {
		fmt.Println("No AI patterns detected in commits.")
		return
	}

	fmt.Printf("Found %d suspicious commits:\n\n", len(warnings))

	for _, w := range warnings {
		fmt.Printf("  %s %s\n", w.Hash, w.Subject)
		for _, reason := range w.Reasons {
			fmt.Printf("    - %s\n", reason)
		}
		fmt.Println()
	}

	fmt.Println("Consider rewriting these commits with: git rebase -i")
}
