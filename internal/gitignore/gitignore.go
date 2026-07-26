package gitignore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type Matcher struct {
	matcher gitignore.Matcher
}

func NewMatcherFromFile(rootDir string) (*Matcher, error) {
	gitignorePath := filepath.Join(rootDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return &Matcher{}, nil
	}

	file, err := os.Open(gitignorePath)
	if err != nil {
		return nil, fmt.Errorf("打开 .gitignore 失败: %w", err)
	}
	defer file.Close()

	var patterns []gitignore.Pattern
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "./")
		pattern := gitignore.ParsePattern(line, nil)
		patterns = append(patterns, pattern)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取 .gitignore 失败: %w", err)
	}

	return &Matcher{
		matcher: gitignore.NewMatcher(patterns),
	}, nil
}

func (m *Matcher) Match(relPath string, isDir bool) bool {
	if m.matcher == nil {
		return false
	}
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	return m.matcher.Match(parts, isDir)
}
