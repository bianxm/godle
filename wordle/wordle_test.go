package wordle

import (
	"testing"

	words "github.com/bianxm/godle/words"
)

func TestNewWordleState(t *testing.T) {
	word := "HELLOTHERE"
	ws := newWordleState(word)
	wordleAsString := string(ws.word[:])

	if wordleAsString != word[:5] {
		t.Errorf("Expected word %s, but got %s", word, wordleAsString)
	}
}

func statusToString(ls letterStatus) string {
	switch ls {
	case none:
		return "none"
	case correct:
		return "correct"
	case present:
		return "present"
	case absent:
		return "absent"
	default:
		return "unknown"
	}
}

func TestNewLetter(t *testing.T) {
	letter := byte('a')
	l := newLetter(letter)
	if l.char != letter {
		t.Errorf("Expecting %v, got %v", letter, l.char)
	}
	if l.status != none {
		t.Errorf("Expecting status none, got %s", statusToString(l.status))
	}
}

func TestNewGuess(t *testing.T) {
	word := "HELPS"
	g := newGuess(word)
	t.Logf("%+v", g)
	t.Logf("New guess: %s", g.string())

	for i, l := range g {
		t.Logf("Letter %d: %c, %s", i, l.char, statusToString(l.status))
		if l.char != word[i] || l.status != none {
			t.Errorf(
				"letter [%d] = %c, %s; want %c, none",
				i,
				l.char,
				statusToString(l.status),
				word[i],
			)
		}
	}
}

func TestUpdateLettersWithWord(t *testing.T) {
	guessWord := "LELOL"
	var word [wordSize]byte
	copy(word[:], "HELLO")
	statuses := []letterStatus{
		present,
		correct,
		correct,
		present,
		absent,
	}

	g := newGuess(guessWord)
	g.updateLettersWithWord(word)

	for i, l := range g {
		// t.Logf(
		// 	"letter[%d] = %c, %s\n",
		// 	i,
		// 	l.char,
		// 	statusToString(l.status),
		// )
		if l.status != statuses[i] {
			t.Errorf(
				"letter [%d] = %c, %s; want %c, %s",
				i,
				l.char,
				statusToString(l.status),
				guessWord[i],
				statusToString(statuses[i]),
			)
		}
	}
}

func TestAppendGuessMaxGuesses(t *testing.T) {
	ws := newWordleState("HELLO")
	for i := 0; i < maxGuesses; i++ {
		word := words.GetWord()
		// word := "LLLLL"
		err := ws.appendGuess(newGuess(word))
		// check currGuess = i+1
		if err != nil {
			t.Errorf(
				"appendGuess() returned error: %s",
				err,
			)
		}
		if ws.currGuess != i+1 {
			t.Errorf(
				"currGuess = %d, want %d",
				ws.currGuess,
				i+1,
			)
		}
		// check ws.guesses[i].string() == word
		if ws.guesses[i].string() != word {
			t.Errorf(
				"appended guess word %s, want %s",
				ws.guesses[i].string(),
				word,
			)
		}
	}
	// add extra one: should fail
	err := ws.appendGuess(newGuess(words.GetWord()))
	// t.Logf("%s", err)
	if err == nil {
		t.Errorf("Should error out for too many guesses, but didn't")
	}
}

func TestAppendGuessError(t *testing.T) {
	ws := newWordleState("HELLO")

	// invalid guess length
	err1 := ws.appendGuess(newGuess("HI"))
	// t.Logf("%s length %d", newGuess("HI").string(), len(newGuess("HI").string()))
	t.Logf("%s", err1)
	if err1 == nil {
		t.Errorf("Request went through, but expecting error 'Invalid guess length'")
	}

	// not a word
	err2 := ws.appendGuess(newGuess("HHHHH"))
	t.Logf("%s", err2)
	if err2 == nil {
		t.Errorf("Request went through, but expecting error 'Invalid word'")
	}
}

func TestIsWordGuessed(t *testing.T) {
	ws := newWordleState("HELLO")
	g := newGuess("HELLO")
	g.updateLettersWithWord(ws.word)
	if err := ws.appendGuess(g); err != nil {
		t.Fatalf("Error: %s", err)
	}
	b := ws.isWordGuessed()
	if !b {
		t.Errorf("Should be true but returned false")
	}
}

func TestShouldEndGameCorrectGuess(t *testing.T) {
	ws := newWordleState("HELLO")
	g := newGuess("HELLO")
	g.updateLettersWithWord(ws.word)
	ws.appendGuess(g)
	if !ws.shouldEndGame() {
		t.Errorf("Should be ending game because correctly guessed")
	}
}

func TestShouldEndGameNoMoreGuesses(t *testing.T) {
	ws := newWordleState("HELLO")
	for i := 0; i < maxGuesses; i++ {
		g := newGuess("YIELD")
		g.updateLettersWithWord(ws.word)
		ws.appendGuess(g)
	}
	t.Logf(
		"Is word guessed:%t\nShould end game: %t",
		ws.isWordGuessed(),
		ws.shouldEndGame(),
	)
	if !ws.shouldEndGame() {
		t.Error("Should end game should be true because no more guesses")
	}
}
