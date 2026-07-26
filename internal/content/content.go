package content

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuan-shuo/copy-paste-ai/assets"
	"github.com/yuan-shuo/copy-paste-ai/internal/config"
)

type FileData struct {
	Path    string
	Content string
	Lang    string
}

type RenderData struct {
	Tree   string
	Files  []FileData
	Prompt string
}

func Build(tree, rootDir string, fileList []string, cfg config.Config) (string, error) {
	var files []FileData
	for _, relPath := range fileList {
		fullPath := filepath.Join(rootDir, relPath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		ext := filepath.Ext(relPath)
		lang := strings.TrimPrefix(ext, ".")
		if lang == "" {
			lang = "text"
		}
		files = append(files, FileData{
			Path:    relPath,
			Content: string(data),
			Lang:    lang,
		})
	}

	tmpl, err := template.New("prompt").Parse(assets.PromptTemplate())
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	renderData := RenderData{
		Tree:   tree,
		Files:  files,
		Prompt: cfg.Prompt.Content,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, renderData); err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
	}

	return buf.String(), nil
}
