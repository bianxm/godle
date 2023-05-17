package main

import (
	"fmt"

	"github.com/bianxm/godle/wordle"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	status := m.renderStatus()
	debug := m.renderDebug()

	game := lipgloss.JoinVertical(lipgloss.Center, status, debug)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, game)
}

const (
	colorPrimary   = lipgloss.Color("#d7dadc")
	colorSecondary = lipgloss.Color("#626262")
	colorSeparator = lipgloss.Color("#9c9c9c")
	colorYellow    = lipgloss.Color("#b59f3b")
	colorGreen     = lipgloss.Color("#538d4e")
)

func statusToColor(ls wordle.LetterStatus) lipgloss.Color {
	switch ls {
	case wordle.None:
		return colorPrimary
	case wordle.Absent:
		return colorSecondary
	case wordle.Present:
		return colorYellow
	case wordle.Correct:
		return colorGreen
	default:
		panic("invalid letter status")
	}
}

func (m *model) renderDebug() string {
	ws := m.ws
	return lipgloss.
		NewStyle().
		Foreground(colorPrimary).
		Render(fmt.Sprintf("[DEBUG] Correct word: %s", string(ws.Word[:])))
}

func (m *model) renderStatus() string {
	return lipgloss.NewStyle().Foreground(colorPrimary).Render(m.status)
}

func renderLetterBox(letter string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(color).
		Foreground(color).
		Render(letter)
}

func renderRowOfBoxes(boxes []string) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, boxes[:]...)
}
