package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

const (
	padding  = 2
	maxWidth = 80
)

var lipglos = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render // nolint:gochecknoglobals

type progressBarModel struct {
	progress progress.Model
	message  string
	max      float64
	ch       chan util.Progress
}

func (m progressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width // - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}

		return m, nil

	case util.Progress:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		m.message = msg.Message
		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(float64(msg.LastStageDur.Milliseconds()) / m.max)

		return m, tea.Batch(cmd, waitForActivity(m.ch)) // continue the cmd

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model) // nolint:forcetypeassert

		return m, cmd

	case stop:
		return m, tea.Quit

	default:
		return m, nil
	}
}

func (m progressBarModel) View() string {
	pad := strings.Repeat(" ", padding)

	return m.message + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + lipglos("Press any key to quit")
}

func RunSingleProgress(p chan util.Progress, max float64) {
	m := progressBarModel{
		progress: progress.New(progress.WithDefaultGradient()),
		max:      max,
		ch:       p,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func (m progressBarModel) Init() tea.Cmd {
	return tea.Batch(waitForActivity(m.ch))
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan util.Progress) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-sub
		if !ok {
			return stop{}
		}

		return msg
	}
}

type stop struct{}
