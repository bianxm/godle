package wordle

import (
	"testing"

	words "github.com/bianxm/godle/words"
)

func TestNewWordleState(t *testing.T) {
	word := "HELLOTHERE"
	ws := NewWordleState(word)
	wordleAsString := string(ws.Word[:])

	if wordleAsString != word[:5] {
		t.Errorf("Expected word %s, but got %s", word, wordleAsString)
	}
	t.Logf("%+v", ws.Alphabet)
}

func statusToString(ls LetterStatus) string {
	switch ls {
	case None:
		return "None"
	case Correct:
		return "Correct"
	case Present:
		return "Present"
	case Absent:
		return "Absent"
	default:
		return "Unknown"
	}
}

func TestNewLetter(t *testing.T) {
	letter := byte('a')
	l := newLetter(letter)
	if l.Char != letter {
		t.Errorf("Expecting %v, got %v", letter, l.Char)
	}
	if l.Status != None {
		t.Errorf("Expecting status none, got %s", statusToString(l.Status))
	}
}

func TestNewGuess(t *testing.T) {
	word := "HELPS"
	g := NewGuess(word)
	t.Logf("%+v", g)
	t.Logf("New guess: %s", g.string())

	for i, l := range g {
		t.Logf("Letter %d: %c, %s", i, l.Char, statusToString(l.Status))
		if l.Char != word[i] || l.Status != None {
			t.Errorf(
				"letter [%d] = %c, %s; want %c, none",
				i,
				l.Char,
				statusToString(l.Status),
				word[i],
			)
		}
	}
}

func TestUpdateLettersWithWord(t *testing.T) {
	guessWord := "LELOL"
	var word [WordSize]byte
	copy(word[:], "HELLO")
	statuses := []LetterStatus{
		Present,
		Correct,
		Correct,
		Present,
		Absent,
	}

	g := NewGuess(guessWord)
	g.UpdateLettersWithWord(word)

	for i, l := range g {
		// t.Logf(
		// 	"letter[%d] = %c, %s\n",
		// 	i,
		// 	l.char,
		// 	statusToString(l.status),
		// )
		if l.Status != statuses[i] {
			t.Errorf(
				"letter [%d] = %c, %s; want %c, %s",
				i,
				l.Char,
				statusToString(l.Status),
				guessWord[i],
				statusToString(statuses[i]),
			)
		}
	}
}

func TestAppendGuessMaxGuesses(t *testing.T) {
	ws := NewWordleState("HELLO")
	for i := 0; i < MaxGuesses; i++ {
		word := words.GetWord()
		// word := "LLLLL"
		err := ws.AppendGuess(NewGuess(word))
		// check currGuess = i+1
		if err != nil {
			t.Errorf(
				"appendGuess() returned error: %s",
				err,
			)
		}
		if ws.CurrGuess != i+1 {
			t.Errorf(
				"currGuess = %d, want %d",
				ws.CurrGuess,
				i+1,
			)
		}
		// check ws.guesses[i].string() == word
		if ws.Guesses[i].string() != word {
			t.Errorf(
				"appended guess word %s, want %s",
				ws.Guesses[i].string(),
				word,
			)
		}
	}
	// add extra one: should fail
	err := ws.AppendGuess(NewGuess(words.GetWord()))
	// t.Logf("%s", err)
	if err == nil {
		t.Errorf("Should error out for too many guesses, but didn't")
	}
}

func TestAppendGuessAlphabetUpdate(t *testing.T) {
	var w [WordSize]byte
	copy(w[:], "HELLO")
	ws := NewWordleState("HELLO")
	word := "HELPS"
	g := NewGuess(word)
	g.UpdateLettersWithWord(w)
	ws.AppendGuess(g)
	t.Logf("%+v", ws.Alphabet)
	statuses := map[byte]LetterStatus{
		'H': Correct,
		'E': Correct,
		'L': Correct,
		'P': Absent,
		'S': Absent,
	}
	for i := 'A'; i <= 'Z'; i++ {
		j := byte(i)
		if ws.Alphabet[j] != statuses[j] {
			t.Errorf(
				"Letter %c: expecting %s, got %s",
				i,
				statusToString(statuses[j]),
				statusToString(ws.Alphabet[j]),
			)
		}
	}
	// check Alphabet
}

func TestAppendGuessError(t *testing.T) {
	ws := NewWordleState("HELLO")

	// invalid guess length
	err1 := ws.AppendGuess(NewGuess("HI"))
	// t.Logf("%s length %d", newGuess("HI").string(), len(newGuess("HI").string()))
	t.Logf("%s", err1)
	if err1 == nil {
		t.Errorf("Request went through, but expecting error 'Invalid guess length'")
	}

	// not a word
	err2 := ws.AppendGuess(NewGuess("HHHHH"))
	t.Logf("%s", err2)
	if err2 == nil {
		t.Errorf("Request went through, but expecting error 'Invalid word'")
	}
}

func TestIsWordGuessed(t *testing.T) {
	ws := NewWordleState("HELLO")
	g := NewGuess("HELLO")
	g.UpdateLettersWithWord(ws.Word)
	if err := ws.AppendGuess(g); err != nil {
		t.Fatalf("Error: %s", err)
	}
	b := ws.IsWordGuessed()
	if !b {
		t.Errorf("Should be true but returned false")
	}
}

func TestShouldEndGameCorrectGuess(t *testing.T) {
	ws := NewWordleState("HELLO")
	g := NewGuess("HELLO")
	g.UpdateLettersWithWord(ws.Word)
	ws.AppendGuess(g)
	t.Logf("%+v", ws.Alphabet)
	if !ws.ShouldEndGame() {
		t.Errorf("Should be ending game because correctly guessed")
	}
}

func TestShouldEndGameNoMoreGuesses(t *testing.T) {
	ws := NewWordleState("HELLO")
	for i := 0; i < MaxGuesses; i++ {
		g := NewGuess("YIELD")
		g.UpdateLettersWithWord(ws.Word)
		ws.AppendGuess(g)
	}
	t.Logf(
		"Is word guessed:%t\nShould end game: %t",
		ws.IsWordGuessed(),
		ws.ShouldEndGame(),
	)
	if !ws.ShouldEndGame() {
		t.Error("Should end game should be true because no more guesses")
	}
}

// func TestInitAlphabet(t *testing.T) {
// 	InitAlphabet()
// 	t.Logf("%+v", Alphabet)
// }
