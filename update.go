// BIANCA !!
package main

import (
	"fmt"

	"github.com/bianxm/godle/wordle"
	"github.com/bianxm/godle/words"
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

		case tea.KeyBackspace:
			m.handleDeleteChar()

		case tea.KeyEnter:
			if m.gameOver {
				// new game initialization
				m.handleResetStatus()
				m.handleResetActiveGuess()
				m.handleResetWordleState()
				m.gameOver = false
				return m, nil
			} else {
				m.handleSubmitActiveGuess()
				m.handleShouldEndGame()
			}

		case tea.KeyRunes:
			if len(msg.Runes) == 1 && !m.gameOver {
				m.handleSubmitChar(msg.Runes[0])
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *model) handleResetWordleState() {
	ws := wordle.NewWordleState(words.GetWord())
	m.ws = &ws
}

func (m *model) handleShouldEndGame() {
	ws := m.ws
	m.gameOver = ws.ShouldEndGame()
	if m.gameOver {
		m.cursor = -1
		if ws.IsWordGuessed() {
			// m.handleSetStatus("Word guessed!\nPress ENTER to restart", 1*time.Second)
			m.handleSetStatus("Word guessed!\nPress ENTER to restart")
		} else {
			// means that there's no more guesses
			// m.handleSetStatus(fmt.Sprintf("No more guesses :( Word was %s\nPress ENTER to restart", string(ws.Word[:])), 1*time.Second)
			m.handleSetStatus(fmt.Sprintf("No more guesses :( Word was %s\nPress ENTER to restart", string(ws.Word[:])))
		}
	}

}

func (m *model) handleSubmitActiveGuess() {
	ws := m.ws
	// only submit until the cursor :)
	g := wordle.NewGuess(string(m.activeGuess[:m.cursor]))
	g.UpdateLettersWithWord(ws.Word)

	err := ws.AppendGuess(g)
	if err != nil {
		// m.handleSetStatus(err.Error(), 1*time.Second)
		m.handleSetStatus(err.Error())
		return
	}
	// fmt.Println(m.ws.Alphabet)
	m.handleResetStatus()
	// reset status to "Guess the word"
	m.handleResetActiveGuess()
}

func (m *model) handleResetActiveGuess() {
	copy(m.activeGuess[:], "")
	m.cursor = 0
}

func (m *model) handleDeleteChar() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) handleSubmitChar(r rune) {
	if m.cursor < wordle.WordSize {
		if 'a' <= r && r <= 'z' {
			r -= 'a' - 'A'
		}
		if 'A' <= r && r <= 'Z' {
			m.activeGuess[m.cursor] = byte(r)
			m.cursor++
		}
	}
}

// handleSetStatus sets the status message, and returns a tea.Cmd that restores the
// default status message after a delay.
// func (m *model) handleSetStatus(msg string, duration time.Duration) tea.Cmd {
func (m *model) handleSetStatus(msg string) {
	m.status = msg
	// if duration > 0 {
	// 	m.statusPending++
	// 	// fmt.Println("hi")
	// 	return tea.Tick(duration, func(time.Time) tea.Msg {
	// 		return msgResetStatus{}
	// 	})
	// }
	// return nil
}

// handleResetStatus immediately resets the status message to its default value.
func (m *model) handleResetStatus() {
	m.status = "Guess the word!"
}

// msgResetStatus is sent when the status line should be reset.
type msgResetStatus struct{}
