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

type optionItem struct{ key, desc string }

func (o optionItem) Title() string       { return o.key }
func (o optionItem) Description() string { return o.desc }
func (o optionItem) FilterValue() string { return o.key }

type selectorModel struct {
	list     list.Model
	selected string
}

// 紧凑委托：仅覆盖 Spacing() 为 0，其余行为沿用默认委托
// 这样能显著缩短每行选项之间的间距

type compactDelegate struct{ list.DefaultDelegate }

func newSelectorModel(options []string, prompt, defaultVal string) selectorModel {
	items := make([]list.Item, 0, len(options))
	for _, k := range options {
		items = append(items, optionItem{key: k, desc: ""})
	}

	// 提升可见性：更高的最小高度，紧凑的行间距
	width := 60
	height := len(items) + 10 // 比原先 +5 更高一些
	if height < 16 {
		height = 16 // 提高最小高度
	}
	if height > 30 {
		height = 30 // 控制上限，避免过高
	}

	delegate := compactDelegate{list.NewDefaultDelegate()}
	delegate.SetSpacing(0)

	l := list.New(items, delegate, width, height)
	l.Title = prompt
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	// 设定默认选中项
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

// fallbackSelectFromStdin 在 Bubble Tea 不可用时，提供简单数字输入回退，并在无默认值时强制有效输入。
func fallbackSelectFromStdin(options []string, prompt, defaultVal string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stdout, "%s\n", prompt)
		for i, opt := range options {
			fmt.Fprintf(os.Stdout, "  [%d] %s\n", i+1, opt)
		}
		if defaultVal != "" {
			fmt.Fprintf(os.Stdout, "请输入编号并回车（默认: %s）: ", defaultVal)
		} else {
			fmt.Fprint(os.Stdout, "请输入编号或名称并回车: ")
		}

		text, err := reader.ReadString('\n')
		if err != nil {
			// 强制交互：读入失败不再静默采用默认，直接报错
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
		// 允许用户直接输入名称
		for _, opt := range options {
			if strings.EqualFold(opt, text) {
				return opt, nil
			}
		}
		fmt.Fprintln(os.Stdout, "无效输入，请重试。")
	}
}

// SelectFromList 使用 Bubble Tea 列表选择一个字符串项。
// 在非 TTY 或设置禁用时直接回退；否则尝试 TUI，失败再回退。
func SelectFromList(options []string, prompt, defaultVal string) (string, error) {
	if len(options) == 0 {
		return "", errors.New("没有可选项")
	}

	// 非交互或显式禁用时，避免先打印 TUI 再回退导致双列表
	if disableTUI() || !isInteractiveTerminal() {
		return fallbackSelectFromStdin(options, prompt, defaultVal)
	}

	m := newSelectorModel(options, prompt, defaultVal)
	p := tea.NewProgram(m, tea.WithOutput(os.Stdout), tea.WithInput(os.Stdin), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err == nil {
		// 注意：Bubble Tea 按值传递 model；Run() 返回最终 model，不能再读外层 m
		var selected string
		switch mm := finalModel.(type) {
		case selectorModel:
			selected = mm.selected
		case *selectorModel:
			selected = mm.selected
		default:
			selected = m.selected // 理论上不该走到这，但给出兜底
		}
		if selected == "" {
			// 在交互环境下，若用户取消（Esc/Ctrl+C），使用回退以获得明确选择
			return fallbackSelectFromStdin(options, prompt, defaultVal)
		}
		return selected, nil
	}
	// Bubble Tea 失败时，使用数字/名称输入回退
	return fallbackSelectFromStdin(options, prompt, defaultVal)
}
