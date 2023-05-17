package wordle

import (
	"errors"

	words "github.com/bianxm/godle/words"
)

const (
	MaxGuesses = 6
	WordSize   = 5
)

type LetterStatus int

const (
	None LetterStatus = iota
	Absent
	Present
	Correct
)

// word - to be guessed
// guesses - array max length 6 of player guesses. made up of:
// guess - string of 6 letters. each letter has state absent/ present/ correct

type WordleState struct {
	Word      [WordSize]byte
	guesses   [MaxGuesses]guess
	currGuess int
}

type guess [WordSize]letter

func (g guess) string() string {
	// var w [wordSize]byte
	str := ""
	for _, l := range g {
		// 	w[i] = l.char
		if 'A' <= l.char && l.char <= 'Z' {
			str += string(l.char)
		}
	}
	// return string(w[:])
	return str
}

type letter struct {
	char   byte
	status LetterStatus
}

func NewWordleState(word string) WordleState {
	w := WordleState{}
	w.Word = [WordSize]byte([]byte(word[:WordSize]))
	return w
}

func newLetter(b byte) letter {
	return letter{char: b}
}

func newGuess(s string) guess {
	// loop over each letter in string
	// convert to letter structs
	var g guess
	for i, l := range s {
		g[i] = newLetter(byte(l))
	}
	return g
}

// GAME LOGIC!
func (g *guess) updateLettersWithWord(word [WordSize]byte) {
	// updates status of the letters in the guess based on a word
	// create a map letter to count
	lc := make(map[byte]int)
	for _, c := range word {
		lc[c] += 1
	}
	// fmt.Println(lc)
	// FIRST iterate through all letters in g and
	// check if word[i] is same as l.char -> correct
	// and subtract from the count map
	for i := range g {
		l := &g[i]
		if word[i] == l.char {
			l.status = Correct
			lc[l.char] -= 1
		}
	}
	// THEN do present/ absent
	for i := range g {
		l := &g[i]
		if l.status != Correct {
			if lc[l.char] > 0 {
				l.status = Present
				lc[l.char] = lc[l.char] - 1
			} else {
				l.status = Absent
			}
		}
	}
}

func (ws *WordleState) appendGuess(g guess) error {
	// return nil if added successfully
	// error if: max guesses already reached, guess isn't long enough, guess isn't valid word
	if ws.currGuess >= MaxGuesses {
		return errors.New("Max guesses reached")
	}

	if len(g.string()) != WordSize {
		return errors.New("Invalid guess length")
	}

	if !words.IsWord(g.string()) {
		return errors.New("Invalid word")
	}

	ws.guesses[ws.currGuess] = g
	ws.currGuess++
	return nil
}

func (ws *WordleState) isWordGuessed() bool {
	// returns true if latest guess is the correct word
	// check ws.guesses[currGuess-1].string() == ws.word
	if ws.currGuess == 0 {
		return false
	}
	if ws.guesses[ws.currGuess-1].string() == string(ws.Word[:]) {
		return true
	}
	return false
}

func (ws *WordleState) shouldEndGame() bool {
	// return true if latest guess is correct
	// or no more guesses are allowed

	if ws.isWordGuessed() || ws.currGuess >= MaxGuesses {
		return true
	}

	return false
}
