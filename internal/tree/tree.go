package tree

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuan-shuo/copy-paste-ai/internal/gitignore"
)

func Generate(rootDir string, gitignoreEnabled bool, m *gitignore.Matcher) string {
	var builder strings.Builder
	builder.WriteString(rootDir)
	builder.WriteString("\n")

	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && d.Name() == ".cpa" {
			return fs.SkipDir
		}

		if d.IsDir() && d.Name() == ".git" {
			return fs.SkipDir
		}

		relPath, _ := filepath.Rel(rootDir, path)
		if relPath == "." {
			return nil
		}

		if gitignoreEnabled {
			if shouldIgnore(relPath, d.IsDir(), m) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		depth := strings.Count(relPath, string(filepath.Separator))
		prefix := strings.Repeat("    ", depth)

		parent := filepath.Dir(relPath)
		entries, err := os.ReadDir(filepath.Join(rootDir, parent))
		if err != nil {
			return nil
		}

		var names []string
		for _, e := range entries {
			if e.Name() == ".cpa" || e.Name() == ".git" {
				continue
			}
			eRelPath, _ := filepath.Rel(rootDir, filepath.Join(parent, e.Name()))
			if gitignoreEnabled {
				if shouldIgnore(eRelPath, e.IsDir(), m) {
					continue
				}
			}
			names = append(names, e.Name())
		}

		isLast := false
		for i, name := range names {
			if name == d.Name() {
				if i == len(names)-1 {
					isLast = true
				}
				break
			}
		}

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if depth == 0 {
			rootEntries, err := os.ReadDir(rootDir)
			if err == nil {
				var rootNames []string
				for _, e := range rootEntries {
					if e.Name() == ".cpa" || e.Name() == ".git" {
						continue
					}
					eRelPath, _ := filepath.Rel(rootDir, filepath.Join(rootDir, e.Name()))
					if gitignoreEnabled {
						if shouldIgnore(eRelPath, e.IsDir(), m) {
							continue
						}
					}
					rootNames = append(rootNames, e.Name())
				}
				for i, name := range rootNames {
					if name == d.Name() && i == len(rootNames)-1 {
						connector = "└── "
					}
				}
			}
		}

		name := d.Name()
		if d.IsDir() {
			name = name + "/"
		}

		builder.WriteString(prefix)
		builder.WriteString(connector)
		builder.WriteString(name)
		builder.WriteString("\n")

		return nil
	})

	return builder.String()
}

func shouldIgnore(relPath string, isDir bool, m *gitignore.Matcher) bool {
	if strings.HasPrefix(relPath, ".git") || strings.HasPrefix(relPath, ".git/") {
		return true
	}
	if m == nil {
		return false
	}
	return m.Match(relPath, isDir)
}
