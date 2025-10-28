package ui

import (
	"strings"
)

type Panel struct {
	Title   string
	Content string
}

func NewPanel(content, title string) Panel {
	return Panel{Title: title, Content: content}
}

func (p Panel) Render() string {
	var b strings.Builder
	if strings.TrimSpace(p.Title) != "" {
		b.WriteString(p.Title)
		b.WriteString("\n")
	}
	b.WriteString(p.Content)
	b.WriteString("\n")
	return b.String()
}