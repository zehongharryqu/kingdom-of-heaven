// the screen where the user picks their room and name
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

type Typewriter struct {
	runes         []rune
	currentText   string // what the user is typing
	counter       int    // frame counter for blink
	confirmedRoom string
	confirmedName string
}

func (t *Typewriter) Update() error {
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	t.runes = ebiten.AppendInputChars(t.runes[:0])
	t.currentText += string(t.runes)

	// Adjust the string to be at most MaxNameChars characters
	if len(t.currentText) > MaxNameChars {
		t.currentText = t.currentText[:MaxNameChars]
	}

	// If the enter key is pressed, confirm the current text
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		if len(t.currentText) > 0 {
			if t.confirmedRoom == "" {
				t.confirmedRoom = t.currentText
				t.currentText = ""
			} else if t.confirmedName == "" {
				t.confirmedName = t.currentText
				t.currentText = ""
			}
		}
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(t.currentText) >= 1 {
			t.currentText = t.currentText[:len(t.currentText)-1]
		}
	}

	t.counter++
	return nil
}

func (t *Typewriter) Draw(screen *ebiten.Image) {
	// Blink the cursor.
	currentTextDisplay := t.currentText
	if t.counter%60 < 30 {
		currentTextDisplay += "_"
	}
	message := "Please enter the name of the room you want to join:\n"
	if t.confirmedRoom != "" {
		message += t.confirmedRoom + "\nPlease enter your name:\n"
	}
	ebitenutil.DebugPrint(screen, message+currentTextDisplay)
}
