package main

import (
	"fmt"

	"github.com/bianxm/godle/wordle"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	status := m.renderStatus()
	grid := m.renderRows()
	debug := m.renderDebug()

	game := lipgloss.JoinVertical(
		lipgloss.Center,
		status,
		grid,
		debug,
	)

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

func (m *model) renderRows() string {
	var rows [wordle.MaxGuesses]string
	ws := m.ws
	for i, g := range ws.Guesses {
		if i < ws.CurrGuess {
			rows[i] = m.renderPastGuess(g)
		} else if i == ws.CurrGuess {
			rows[i] = m.renderActiveGuess()
		} else {
			rows[i] = m.renderFutureGuess()
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows[:]...)
}

func (m *model) renderPastGuess(g wordle.Guess) string {
	var letterBoxes [wordle.WordSize]string
	for i, l := range g {
		letterBoxes[i] = renderLetterBox(string(l.Char), statusToColor(l.Status))
	}
	return renderRowOfBoxes(letterBoxes[:])
}

func (m *model) renderFutureGuess() string {
	var letterBoxes [wordle.WordSize]string
	for i := 0; i < wordle.WordSize; i++ {
		letterBoxes[i] = renderLetterBox(" ", colorPrimary)
	}
	return renderRowOfBoxes(letterBoxes[:])
}

func (m *model) renderActiveGuess() string {
	var letterBoxes [wordle.WordSize]string
	for i, char := range m.activeGuess {
		var letter string
		if i < m.cursor {
			letter = string(char)
		} else if i == m.cursor {
			letter = "_"
		} else {
			letter = " "
		}

		letterBoxes[i] = renderLetterBox(letter, colorPrimary)
	}

	return renderRowOfBoxes(letterBoxes[:])
}
