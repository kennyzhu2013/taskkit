package ui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type optionItem struct{ key, name, desc string }

type compactDelegate struct{ list.DefaultDelegate }

type selectorModel struct {
	list     list.Model
	selected string
}

func (o optionItem) Title() string       { return o.name }
func (o optionItem) Description() string { return o.desc }
func (o optionItem) FilterValue() string { return o.key }

// 自定义显示名版本：Title 使用映射中的显示名，返回键名
func newSelectorModelWithDisplay(options []string, prompt, defaultVal string, displayNames map[string]string) selectorModel {
	items := make([]list.Item, 0, len(options))
	for _, k := range options {
		name := k
		if v, ok := displayNames[k]; ok && strings.TrimSpace(v) != "" {
			name = v
		}
		items = append(items, optionItem{key: k, name: name, desc: ""})
	}

	width := 60
	height := len(items) + 10
	if height < 16 {
		height = 16
	}
	if height > 30 {
		height = 30
	}

	delegate := compactDelegate{list.NewDefaultDelegate()}
	delegate.SetSpacing(0)

	l := list.New(items, delegate, width, height)
	l.Title = prompt
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowPagination(false)

	for i, it := range items {
		if oi, ok := it.(optionItem); ok && oi.key == defaultVal {
			l.Select(i)
			break
		}
	}
	return selectorModel{list: l}
}

func (m selectorModel) Init() tea.Cmd { return nil }

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if it, ok := m.list.SelectedItem().(optionItem); ok {
				m.selected = it.key
			}
			return m, tea.Quit
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectorModel) View() string { return m.list.View() }

// 带显示名的回退：列表显示为 显示名(key)，输入既可是编号、键名也可是显示名
func fallbackSelectFromStdinWithDisplay(options []string, prompt, defaultVal string, displayNames map[string]string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stdout, "%s\n", prompt)
		for i, k := range options {
			name := k
			if v, ok := displayNames[k]; ok && strings.TrimSpace(v) != "" {
				name = v
			}
			fmt.Fprintf(os.Stdout, "  [%d] %s (%s)\n", i+1, name, k)
		}
		if defaultVal != "" {
			fmt.Fprintf(os.Stdout, "请输入编号并回车（默认: %s）: ", defaultVal)
		} else {
			fmt.Fprint(os.Stdout, "请输入编号或名称并回车: ")
		}

		text, err := reader.ReadString('\n')
		if err != nil {
			return "", errors.New("需要交互选择但未读取到输入；请使用命令行参数指定或在交互终端运行")
		}
		text = strings.TrimSpace(text)
		if text == "" {
			if defaultVal != "" {
				return defaultVal, nil
			}
			fmt.Fprintln(os.Stdout, "无默认项，请输入有效编号或名称。")
			continue
		}
		if idx, err := strconv.Atoi(text); err == nil {
			idx--
			if idx >= 0 && idx < len(options) {
				return options[idx], nil
			}
			fmt.Fprintln(os.Stdout, "编号超出范围，请重试。")
			continue
		}
		// 键名或显示名均可匹配
		for _, k := range options {
			if strings.EqualFold(k, text) {
				return k, nil
			}
			if v, ok := displayNames[k]; ok && strings.EqualFold(strings.TrimSpace(v), text) {
				return k, nil
			}
		}
		fmt.Fprintln(os.Stdout, "无效输入，请重试。")
	}
}

// SelectFromListWithDisplay 支持用显示名渲染列表，同时返回键名
func SelectFromListWithDisplay(options []string, prompt, defaultVal string, displayNames map[string]string) (string, error) {
	if len(options) == 0 {
		return "", errors.New("没有可选项")
	}
	if disableTUI() || !isInteractiveTerminal() {
		return fallbackSelectFromStdinWithDisplay(options, prompt, defaultVal, displayNames)
	}
	m := newSelectorModelWithDisplay(options, prompt, defaultVal, displayNames)
	// 移除 AltScreen，避免 Windows 下切换至回退输入后键盘失效
	p := tea.NewProgram(m, tea.WithOutput(os.Stdout), tea.WithInput(os.Stdin))
	finalModel, err := p.Run()
	if err == nil {
		var selected string
		switch mm := finalModel.(type) {
		case selectorModel:
			selected = mm.selected
		case *selectorModel:
			selected = mm.selected
		}
		if selected == "" {
			return fallbackSelectFromStdinWithDisplay(options, prompt, defaultVal, displayNames)
		}
		return selected, nil
	}
	return fallbackSelectFromStdinWithDisplay(options, prompt, defaultVal, displayNames)
}
