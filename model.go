package main

import (
	"github.com/bianxm/godle/wordle"
	"github.com/bianxm/godle/words"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	ws *wordle.WordleState

	activeGuess [wordle.WordSize]byte
	cursor      int

	status        string
	statusPending int

	width  int
	height int

	gameOver bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	ws := wordle.NewWordleState(words.GetWord())
	m := model{
		ws:     &ws,
		status: "Guess the word!",
	}
	return m
}
