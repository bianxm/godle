package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgResetStatus:
		// If there is more than one pending status message, that means
		// something else is currently displaying a status message, so we don't
		// want to overwrite it.
		m.statusPending--
		if m.statusPending == 0 {
			m.handleResetStatus()
		}

	// Handle keypresses
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlD:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// handleSetStatus sets the status message, and returns a tea.Cmd that restores the
// default status message after a delay.
func (m *model) handleSetStatus(msg string, duration time.Duration) tea.Cmd {
	m.status = msg
	if duration > 0 {
		m.statusPending++
		return tea.Tick(duration, func(time.Time) tea.Msg {
			return msgResetStatus{}
		})
	}
	return nil
}

// handleResetStatus immediately resets the status message to its default value.
func (m *model) handleResetStatus() {
	m.status = "Guess the word!"
}

// msgResetStatus is sent when the status line should be reset.
type msgResetStatus struct{}
