package walker

import (
	"os"
	"path/filepath"
	"strings"
)

// SupportedExtensions defines the file extensions we process
var SupportedExtensions = map[string]bool{
	".ts":  true,
	".tsx": true,
	".js":  true,
	".jsx": true,
	".py":  true,
}

// IgnorePatterns defines directories to skip
var IgnorePatterns = []string{
	".git",
	"node_modules",
	"__pycache__",
	".venv",
	"venv",
	"dist",
	"build",
	".next",
	"coverage",
	".cache",
}

// FileInfo holds information about a file to process
type FileInfo struct {
	Path string
	Ext  string
}

// Walk traverses the given path and returns all supported files
func Walk(root string) ([]FileInfo, error) {
	var files []FileInfo

	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	// Single file
	if !info.IsDir() {
		ext := strings.ToLower(filepath.Ext(root))
		if SupportedExtensions[ext] {
			files = append(files, FileInfo{Path: root, Ext: ext})
		}
		return files, nil
	}

	// Directory
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		if info.IsDir() {
			name := info.Name()
			for _, pattern := range IgnorePatterns {
				if name == pattern {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check extension
		ext := strings.ToLower(filepath.Ext(path))
		if SupportedExtensions[ext] {
			files = append(files, FileInfo{Path: path, Ext: ext})
		}

		return nil
	})

	return files, err
}

// IsJavaScript returns true if the extension is JS/TS
func (f FileInfo) IsJavaScript() bool {
	return f.Ext == ".js" || f.Ext == ".jsx" || f.Ext == ".ts" || f.Ext == ".tsx"
}

// IsPython returns true if the extension is Python
func (f FileInfo) IsPython() bool {
	return f.Ext == ".py"
}
