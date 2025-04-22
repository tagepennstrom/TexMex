// OBS måste ladda ner: go get github.com/eiannone/keyboard

package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

func processInput() {
	// Open the keyboard for reading key events.
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	doc := NewDocument()

	fmt.Println("Press keys. (Press '0' to exit)")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		// Quit på '0'
		if char == '0' {
			fmt.Println("Exiting...")
			break
		}
		if char == '9' {
			fmt.Print("\033[H\033[2J")
			doc.DisplayWithCursor()

			doc.Debug(true)
		} else if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			doc.letterAction(char)
		} else if key == keyboard.KeyArrowLeft {
			doc.leftArrowAction()
		} else if key == keyboard.KeyArrowRight {
			doc.rightArrowAction()
		} else if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {
			doc.deleteAction()
		} else if key == keyboard.KeySpace {
			doc.letterAction(' ')
		}
	}
}

func (d *Document) letterAction(letter rune) {
	fmt.Print("\033[H\033[2J")
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	randomNumber := r.Intn(10000) + 1

	d.Insert(string(letter), randomNumber)
	d.DisplayWithCursor()
}

func (d *Document) leftArrowAction() {
	fmt.Print("\033[H\033[2J")

	d.CursorBackwards()
	d.DisplayWithCursor()

}

func (d *Document) rightArrowAction() {
	fmt.Print("\033[H\033[2J")

	d.CursorForward()
	d.DisplayWithCursor()

}

func (d *Document) deleteAction() {
	fmt.Print("\033[H\033[2J")

	d.Delete()
	d.DisplayWithCursor()

}
