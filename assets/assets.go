package assets

import "embed"

//go:embed skel
var skelFS embed.FS

func ConfigTemplate() string {
	data, err := skelFS.ReadFile("skel/config.toml")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func PromptTemplate() string {
	data, err := skelFS.ReadFile("skel/prompt.tmpl")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func DefaultPrompt() string {
	data, err := skelFS.ReadFile("skel/prompt_default.md")
	if err != nil {
		panic(err)
	}
	return string(data)
}
