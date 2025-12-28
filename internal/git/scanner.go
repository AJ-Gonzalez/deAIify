package git

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// CommitWarning describes an issue with a commit
type CommitWarning struct {
	Hash    string
	Subject string
	Reasons []string
}

// AI footers to detect
var aiFooters = []string{
	"generated with",
	"co-authored-by: claude",
	"co-authored-by: chatgpt",
	"co-authored-by: copilot",
	"co-authored-by: gemini",
	"co-authored-by: gpt",
	"created by ai",
	"written by ai",
}

// AI-like commit message patterns
var aiCommitPatterns = []string{
	"this commit",
	"this change",
	"this update",
	"this patch",
	"here we",
	"this pr",
	"this pull request",
	"in this commit",
}

// Emoji regex
var emojiRe = regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F1E0}-\x{1F1FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)

// Common commit emojis (text-based)
var textEmojis = []string{
	":sparkles:", ":bug:", ":memo:", ":rocket:", ":art:",
	":fire:", ":zap:", ":lipstick:", ":tada:", ":white_check_mark:",
	":lock:", ":bookmark:", ":rotating_light:", ":construction:",
	":green_heart:", ":arrow_down:", ":arrow_up:", ":pushpin:",
	":recycle:", ":heavy_plus_sign:", ":heavy_minus_sign:",
	":wrench:", ":hammer:", ":globe_with_meridians:", ":pencil2:",
	":poop:", ":rewind:", ":twisted_rightwards_arrows:", ":package:",
	":truck:", ":page_facing_up:", ":boom:", ":bento:",
	":wheelchair:", ":bulb:", ":beers:", ":speech_balloon:",
	":card_file_box:", ":loud_sound:", ":mute:", ":busts_in_silhouette:",
}

// IsGitRepo checks if the current directory is a git repository
func IsGitRepo(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// ScanCommits scans recent commits for AI patterns
func ScanCommits(path string, count int) ([]CommitWarning, error) {
	if count <= 0 {
		count = 20
	}

	// Get recent commits with full message
	cmd := exec.Command("git", "-C", path, "log",
		"--format=%H%n%s%n%b%n---COMMIT_END---",
		"-n", strconv.Itoa(count))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseCommits(string(output)), nil
}

// parseCommits parses git log output and checks for AI patterns
func parseCommits(output string) []CommitWarning {
	var warnings []CommitWarning

	commits := strings.Split(output, "---COMMIT_END---")
	for _, commit := range commits {
		commit = strings.TrimSpace(commit)
		if commit == "" {
			continue
		}

		lines := strings.SplitN(commit, "\n", 3)
		if len(lines) < 2 {
			continue
		}

		hash := lines[0]
		subject := lines[1]
		body := ""
		if len(lines) > 2 {
			body = lines[2]
		}

		fullMessage := subject + "\n" + body
		lowerMessage := strings.ToLower(fullMessage)

		var reasons []string

		// Check for emoji
		if emojiRe.MatchString(fullMessage) {
			reasons = append(reasons, "contains emoji")
		}

		// Check for text-based emoji
		for _, emoji := range textEmojis {
			if strings.Contains(lowerMessage, emoji) {
				reasons = append(reasons, "contains gitmoji ("+emoji+")")
				break
			}
		}

		// Check for AI footers
		for _, footer := range aiFooters {
			if strings.Contains(lowerMessage, footer) {
				reasons = append(reasons, "contains AI footer: "+footer)
			}
		}

		// Check for AI-like patterns
		for _, pattern := range aiCommitPatterns {
			if strings.Contains(lowerMessage, pattern) {
				reasons = append(reasons, "uses AI-like phrasing: \""+pattern+"\"")
				break
			}
		}

		if len(reasons) > 0 {
			// Truncate hash for display
			shortHash := hash
			if len(hash) > 7 {
				shortHash = hash[:7]
			}
			warnings = append(warnings, CommitWarning{
				Hash:    shortHash,
				Subject: truncate(subject, 60),
				Reasons: reasons,
			})
		}
	}

	return warnings
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
