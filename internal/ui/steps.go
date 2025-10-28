package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type StepStatus int

const (
	StatusPending StepStatus = iota
	StatusRunning
	StatusDone
	StatusError
	StatusSkipped
)

type StepItem struct {
	Key     string
	Label   string
	Detail  string
	Status  StepStatus
	Percent float64
}

type startStepMsg struct{ key, detail string }
type completeStepMsg struct{ key, detail string }
type errorStepMsg struct{ key, detail string }
type skipStepMsg struct{ key, detail string }
type progressMsg struct {
	key     string
	percent float64
}
type quitMsg struct{}

type stepModel struct {
	title string
	order []string
	steps map[string]*StepItem
	prog  progress.Model
	done  bool
}

func NewStepModel(title string, items []StepItem) stepModel {
	m := stepModel{
		title: title,
		steps: make(map[string]*StepItem, len(items)),
		prog:  progress.New(),
	}
	for i := range items {
		item := items[i]
		m.order = append(m.order, item.Key)
		copy := item
		m.steps[item.Key] = &copy
	}
	return m
}

func (m stepModel) Init() tea.Cmd { return nil }

func (m stepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startStepMsg:
		if s := m.steps[msg.key]; s != nil {
			s.Status = StatusRunning
			s.Detail = msg.detail
			s.Percent = 0
		}
	case completeStepMsg:
		if s := m.steps[msg.key]; s != nil {
			s.Status = StatusDone
			s.Detail = msg.detail
			s.Percent = 1
		}
		m.checkDone()
	case errorStepMsg:
		if s := m.steps[msg.key]; s != nil {
			s.Status = StatusError
			s.Detail = msg.detail
		}
		m.checkDone()
	case skipStepMsg:
		if s := m.steps[msg.key]; s != nil {
			s.Status = StatusSkipped
			s.Detail = msg.detail
		}
	case progressMsg:
		if s := m.steps[msg.key]; s != nil {
			s.Percent = msg.percent
		}
	case quitMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *stepModel) checkDone() {
	for _, k := range m.order {
		s := m.steps[k]
		if s.Status != StatusDone && s.Status != StatusSkipped {
			return
		}
	}
	m.done = true
}

func (m stepModel) View() string {
	var out string
	if m.title != "" {
		out += m.title + "\n\n"
	}
	for _, k := range m.order {
		s := m.steps[k]
		status := map[StepStatus]string{
			StatusPending: "[ ]",
			StatusRunning: "[>]",
			StatusDone:    "[✓]",
			StatusError:   "[x]",
			StatusSkipped: "[-]",
		}[s.Status]
		line := fmt.Sprintf("%s %s", status, s.Label)
		if s.Detail != "" {
			line += " — " + s.Detail
		}
		out += line + "\n"
		if s.Status == StatusRunning && s.Percent > 0 && s.Percent < 1 {
			out += m.prog.ViewAs(s.Percent) + "\n"
		}
	}
	return out
}

// StepUI wraps a tea.Program and exposes helper methods to send messages.
type StepUI struct {
	Program *tea.Program
}

// StartStepUI launches the TUI program asynchronously and returns a controller.
func StartStepUI(title string, items []StepItem) (*StepUI, error) {
	m := NewStepModel(title, items)
	p := tea.NewProgram(m)
	go func() { _, _ = p.Run() }()
	return &StepUI{Program: p}, nil
}

func (u *StepUI) Start(key, detail string) { u.Program.Send(startStepMsg{key: key, detail: detail}) }
func (u *StepUI) Complete(key, detail string) {
	u.Program.Send(completeStepMsg{key: key, detail: detail})
}
func (u *StepUI) Error(key, detail string) { u.Program.Send(errorStepMsg{key: key, detail: detail}) }
func (u *StepUI) Skip(key, detail string)  { u.Program.Send(skipStepMsg{key: key, detail: detail}) }
func (u *StepUI) ProgressPercent(key string, percent float64) {
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}
	u.Program.Send(progressMsg{key: key, percent: percent})
}
func (u *StepUI) ProgressBytes(key string, downloaded, total int64) {
	if total <= 0 {
		return
	}
	u.ProgressPercent(key, float64(downloaded)/float64(total))
}
func (u *StepUI) Stop() { u.Program.Send(quitMsg{}) }
