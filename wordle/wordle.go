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
// guess - string of 6 letters
// letter has state absent/ present/ correct

type WordleState struct {
	Word      [WordSize]byte
	Guesses   [MaxGuesses]Guess
	CurrGuess int
}

type Guess [WordSize]letter

func (g Guess) string() string {
	// var w [wordSize]byte
	str := ""
	for _, l := range g {
		// 	w[i] = l.char
		if 'A' <= l.Char && l.Char <= 'Z' {
			str += string(l.Char)
		}
	}
	// return string(w[:])
	return str
}

type letter struct {
	Char   byte
	Status LetterStatus
}

func NewWordleState(word string) WordleState {
	w := WordleState{}
	w.Word = [WordSize]byte([]byte(word[:WordSize]))
	return w
}

func newLetter(b byte) letter {
	return letter{Char: b}
}

func NewGuess(s string) Guess {
	// loop over each letter in string
	// convert to letter structs
	var g Guess
	for i, l := range s {
		g[i] = newLetter(byte(l))
	}
	return g
}

// GAME LOGIC!
func (g *Guess) UpdateLettersWithWord(word [WordSize]byte) {
	// updates status of the letters in the guess based on a word
	// create a map of letter to count
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
		if word[i] == l.Char {
			l.Status = Correct
			lc[l.Char] -= 1
		}
	}
	// THEN do present/ absent
	for i := range g {
		l := &g[i]
		if l.Status != Correct {
			if lc[l.Char] > 0 {
				l.Status = Present
				lc[l.Char] = lc[l.Char] - 1
			} else {
				l.Status = Absent
			}
		}
	}
}

func (ws *WordleState) AppendGuess(g Guess) error {
	// return nil if added successfully
	// error if: max guesses already reached, guess isn't long enough, guess isn't valid word
	if ws.CurrGuess >= MaxGuesses {
		return errors.New("Max guesses reached")
	}

	if len(g.string()) != WordSize {
		return errors.New("Invalid guess length")
	}

	if !words.IsWord(g.string()) {
		return errors.New("Invalid word")
	}

	ws.Guesses[ws.CurrGuess] = g
	ws.CurrGuess++
	return nil
}

func (ws *WordleState) IsWordGuessed() bool {
	// returns true if latest guess is the correct word
	// check ws.guesses[currGuess-1].string() == ws.word
	if ws.CurrGuess == 0 {
		return false
	}
	return ws.Guesses[ws.CurrGuess-1].string() == string(ws.Word[:])
}

func (ws *WordleState) ShouldEndGame() bool {
	// return true if latest guess is correct
	// or no more guesses are allowed

	return ws.IsWordGuessed() || ws.CurrGuess >= MaxGuesses
}
