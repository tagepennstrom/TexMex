// OBS måste ladda ner: go get github.com/eiannone/keyboard

package crdt

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

func liveUserDemo(d *Document) {
	arr := []string{"(Hello Antwan!) ",
					"(Hungry Hippos) ",
					"(We the best music) ",
					"(God did) ",
					"(Another one) ",
					"(Bless up) ",
					"(Major key) ",
					"(Stay away from 'they') ",
					"(Secure the bag) ",
					"(Don't play yourself) ",
					"(They don't want you to win) ",
					"(Lion!) ",
					"(I'm up to something) ",
					"(You smart) ",
					"(You loyal) ",
					"(Celebrate success) ",
					"(Keep going) ",}

	wordIndexChosen := rand.Int() % len(arr)

	word := arr[wordIndexChosen]
	wordLen := len(word)

	pos := rand.Int() % d.Textcontent.Length 

	r2 := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < wordLen; i++ {

		d.LoadInsert(string(word[i]), pos+i, 3)
		fmt.Print("\033[H\033[2J")
		d.DisplayWithCursor()

		time.Sleep(time.Duration(r2.Intn(500)) * time.Millisecond)
	}

}

func processInput() {
	// Open the keyboard for reading key events.
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	doc := NewDocument()

	var debugOn bool = false

	for {
		fmt.Println("(Press '0' to exit. Press '9' for debug screen. Press '1' to simulate another user typing, '2' to reset cursor, '3' to reset coordinates)")

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
			debugOn = !debugOn
			fmt.Print("\033[H\033[2J")
			doc.DisplayWithCursor()

		} else if char == '1' {
			go liveUserDemo(&doc)

		} else if char == '2' {
				doc.MoveCursor(0)
				fmt.Print("\033[H\033[2J")
				doc.DisplayWithCursor()
		
		} else if char == '3' {
			doc.CordReset()

		} else if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			doc.letterAction(char)
		} else if key == keyboard.KeyArrowLeft {
			doc.leftArrowAction()
			//fmt.Print("This is the letter to my left: ", doc.CursorPosition.Prev.Letter, "\n")
		} else if key == keyboard.KeyArrowRight {
			doc.rightArrowAction()
			//fmt.Print("This is the letter to my right: ", doc.CursorPosition.Next.Letter, "\n")
		} else if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {
			doc.deleteAction()
		} else if key == keyboard.KeySpace {
			doc.letterAction(' ')
		}

		if debugOn {
			doc.Debug(true)
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
