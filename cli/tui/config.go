package tui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
)

// nolint
var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type sshConfigModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

const Enter = "enter"

func initialModel() sshConfigModel {
	m := sshConfigModel{
		inputs: make([]textinput.Model, len(inputKeys)),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Address"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Username"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 3:
			//TODO: file picker
			t.Placeholder = "SSH Key"
		}

		m.inputs[i] = t
	}

	return m
}

var inputKeys = []string{"addr", "user", "pass", "key"} // nolint

func (m sshConfigModel) Init() tea.Cmd {
	if addr, ok := config.Get("ssh", "addr"); ok {
		m.inputs[0].SetValue(addr)
	}

	if user, ok := config.Get("ssh", "user"); ok {
		m.inputs[1].SetValue(user)
	}

	if pass, ok := config.Get("ssh", "pass"); ok {
		m.inputs[2].SetValue(pass)
	}

	if key, ok := config.Get("ssh", "key"); ok {
		m.inputs[3].SetValue(key)
	}

	return textinput.Blink
}

func (m sshConfigModel) Update(rmsg tea.Msg) (tea.Model, tea.Cmd) {
	msg, ok := rmsg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "esc", "q":
		return m, tea.Quit

	// Change cursor mode
	case "ctrl+r":
		m.cursorMode++
		if m.cursorMode > cursor.CursorHide {
			m.cursorMode = cursor.CursorBlink
		}

		cmds := make([]tea.Cmd, len(m.inputs))

		for i := range m.inputs {
			cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
		}

		return m, tea.Batch(cmds...)

	// Set focus to next input
	case "tab", "shift+tab", Enter, "up", "down":
		return m.SwitchInputs(msg.String())
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func SaveData(inputs []textinput.Model) error {
	for i := range inputs {
		SaveInput(i, inputs[i].Value())
	}

	return config.Write()
}

func SaveInput(i int, v string) {
	if i >= len(inputKeys) {
		return
	}

	if v == "" {
		config.Unset("ssh", inputKeys[i])
	} else {
		config.Set("ssh", inputKeys[i], v)
	}
}

func (m sshConfigModel) SwitchInputs(s string) (tea.Model, tea.Cmd) {
	// Did the user press enter while the submit button was focused?
	// If so, exit.
	if s == Enter && m.focusIndex == len(m.inputs) {
		if err := SaveData(m.inputs); err != nil {
			fmt.Println(err)
		}

		return m, tea.Quit
	}

	// Cycle indexes
	if s == "up" || s == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	cmds := make([]tea.Cmd, len(m.inputs))

	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle

			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}

	return m, tea.Batch(cmds...)
}
func (m *sshConfigModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m sshConfigModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())

		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func ConfigSSH() error {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		return fmt.Errorf("could not start program: %w", err)
	}

	return nil
}
