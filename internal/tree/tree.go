package tree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuan-shuo/copy-paste-ai/internal/gitignore"
)

func Generate(rootDir string, gitignoreEnabled bool, m *gitignore.Matcher) string {
	var builder strings.Builder
	builder.WriteString(rootDir)
	builder.WriteString("\n")

	entries, _ := readSortedEntries(rootDir, rootDir, gitignoreEnabled, m)
	for i, entry := range entries {
		isLast := i == len(entries)-1
		writeEntry(&builder, rootDir, rootDir, entry, "", isLast, gitignoreEnabled, m)
	}

	return builder.String()
}

type treeEntry struct {
	name  string
	isDir bool
}

func writeEntry(builder *strings.Builder, rootDir, currentDir string, entry treeEntry, prefix string, isLast bool, gitignoreEnabled bool, m *gitignore.Matcher) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	}

	builder.WriteString(prefix)
	builder.WriteString(connector)
	builder.WriteString(entry.name)
	if entry.isDir {
		builder.WriteString("/")
	}
	builder.WriteString("\n")

	if entry.isDir {
		childDir := filepath.Join(currentDir, entry.name)
		entries, _ := readSortedEntries(childDir, rootDir, gitignoreEnabled, m)
		for i, child := range entries {
			childIsLast := i == len(entries)-1
			writeEntry(builder, rootDir, childDir, child, childPrefix, childIsLast, gitignoreEnabled, m)
		}
	}
}

func readSortedEntries(dir, rootDir string, gitignoreEnabled bool, m *gitignore.Matcher) ([]treeEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []treeEntry
	var files []treeEntry

	for _, e := range entries {
		if e.Name() == ".cpa" || e.Name() == ".git" {
			continue
		}
		eRelPath, _ := filepath.Rel(rootDir, filepath.Join(dir, e.Name()))
		if gitignoreEnabled && shouldIgnore(eRelPath, e.IsDir(), m) {
			continue
		}
		entry := treeEntry{name: e.Name(), isDir: e.IsDir()}
		if e.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].name) < strings.ToLower(dirs[j].name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].name) < strings.ToLower(files[j].name)
	})

	return append(dirs, files...), nil
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
