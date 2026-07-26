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

		if d.IsDir() && (d.Name() == ".cpa" || d.Name() == ".git") {
			return fs.SkipDir
		}

		relPath, _ := filepath.Rel(rootDir, path)
		if relPath == "." {
			return nil
		}

		if gitignoreEnabled && shouldIgnore(relPath, d.IsDir(), m) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		depth := strings.Count(relPath, string(filepath.Separator))
		prefix := strings.Repeat("    ", depth)
		connector := findConnector(rootDir, filepath.Dir(relPath), d.Name(), gitignoreEnabled, m)

		name := d.Name()
		if d.IsDir() {
			name += "/"
		}

		builder.WriteString(prefix)
		builder.WriteString(connector)
		builder.WriteString(name)
		builder.WriteString("\n")

		return nil
	})

	return builder.String()
}

func listFilteredDir(dir, rootDir string, gitignoreEnabled bool, m *gitignore.Matcher) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.Name() == ".cpa" || e.Name() == ".git" {
			continue
		}
		eRelPath, _ := filepath.Rel(rootDir, filepath.Join(dir, e.Name()))
		if gitignoreEnabled && shouldIgnore(eRelPath, e.IsDir(), m) {
			continue
		}
		names = append(names, e.Name())
	}
	return names, nil
}

func findConnector(rootDir, parentDir, name string, gitignoreEnabled bool, m *gitignore.Matcher) string {
	names, err := listFilteredDir(filepath.Join(rootDir, parentDir), rootDir, gitignoreEnabled, m)
	if err != nil {
		return "├── "
	}
	for i, n := range names {
		if n == name {
			if i == len(names)-1 {
				return "└── "
			}
			break
		}
	}
	return "├── "
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
