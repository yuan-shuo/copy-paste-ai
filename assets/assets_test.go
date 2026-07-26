package assets

import (
	"strings"
	"testing"
)

func TestConfigTemplate(t *testing.T) {
	content := ConfigTemplate()
	if content == "" {
		t.Fatal("ConfigTemplate() 返回空字符串")
	}
	if !strings.Contains(content, "Copy-Paste AI 配置文件") {
		t.Error("配置模板缺少标题注释")
	}
	if !strings.Contains(content, "[default]") {
		t.Error("配置模板缺少 [default] 段")
	}
	if !strings.Contains(content, "[file_aliases]") {
		t.Error("配置模板缺少 [file_aliases] 段")
	}
	if !strings.Contains(content, "[prompt]") {
		t.Error("配置模板缺少 [prompt] 段")
	}
	if !strings.Contains(content, "[gitignore]") {
		t.Error("配置模板缺少 [gitignore] 段")
	}
}

func TestPromptTemplate(t *testing.T) {
	content := PromptTemplate()
	if content == "" {
		t.Fatal("PromptTemplate() 返回空字符串")
	}
	if !strings.Contains(content, "{{.Tree}}") {
		t.Error("提示词模板缺少 {{.Tree}} 占位符")
	}
	if !strings.Contains(content, "{{- range .Files}}") {
		t.Error("提示词模板缺少 {{- range .Files}} 循环")
	}
	if !strings.Contains(content, "{{.Prompt}}") {
		t.Error("提示词模板缺少 {{.Prompt}} 占位符")
	}
	if !strings.Contains(content, "```") {
		t.Error("提示词模板缺少代码块标记")
	}
}

func TestDefaultPrompt(t *testing.T) {
	content := DefaultPrompt()
	if content == "" {
		t.Fatal("DefaultPrompt() 返回空字符串")
	}
	if !strings.Contains(content, "提示") {
		t.Error("默认提示词缺少标题")
	}
	if !strings.Contains(content, "请根据以上项目文件树和代码") {
		t.Error("默认提示词缺少核心指令")
	}
	if !strings.Contains(content, "完整的代码") {
		t.Error("默认提示词缺少代码要求")
	}
}
