package main

import (
	"fmt"
	"os"
)

type Document struct {
	CursorPosition *Item
	Textcontent    LinkedList
}

type LinkedList struct {
	Head   *Item
	Tail   *Item
	Length int
}

type Item struct {
	Letter     string
	Coordinate []int
	ID         int
	Prev       *Item
	Next       *Item
}

func (ll *LinkedList) Append(newItem *Item) {
	if ll.Tail == nil {
		ll.Head = newItem
		ll.Tail = newItem
		return
	}
	ll.Tail.Next = newItem
	newItem.Prev = ll.Tail
	ll.Tail = newItem
}

func CompareIndexes(coord1 []int, coord2 []int) bool {

	len1 := len(coord1)
	len2 := len(coord2)
	count := 0

	for count < len1 && count < len2 {
		//Var tvungen att ändra på true och false här. Vi vill ju att denna funktion skall hålla på tills C1 är större, men tror fortfarande
		//det finns något underliggande logisk fel i denna men kan ej komma på nu. Nu hamnar de iallafall rätt i ll om man insertar längst bak.
		//Inser nu att det uppstår ett problem vid deltion nu med denna implementation, Den tolkar allt som en tiebreaker. 
		//Om man har false när c1 är störst och true när c2 är minst verkar inget
		if coord1[count] > coord2[count] {
			return true //c1 biggest

		//La till "<=" för detta verkar funka för tillfället, kanske ta bort sen? eller så ska det vara såhär, förmodligen inte
		} else if coord1[count] <= coord2[count] {
			return false //c2 biggest
		}
		count++

	}
	if len1 == len2 {
		fmt.Errorf("Coordinates have the same size")
		println("Error: Coordinates can't have the same size. OBS Måste implementera tiebreaker")
		os.Exit(1)
	}
	return len1 < len2
}

func getCoordinateProperties(prevCord []int, nextCord []int) []int {

	prevLen := len(prevCord)
	nextLen := len(nextCord)
	count := 0
	newCoordinate := []int{}

	for count < prevLen && count < nextLen {
		comparison := nextCord[count] - prevCord[count]

		if comparison == 0 {
			newCoordinate = append(newCoordinate, prevCord[count])

			// Case 1. Plats finns
		} else if comparison > 1 {
			newCoordinate = append(newCoordinate, prevCord[count]+1)

			//Case Z
			if count+1 < prevLen {
				newCoordinate = append(newCoordinate, prevCord[count+1]+1)

			}
			return newCoordinate

			// Case 2 and Case X. Ingen plats finns
		} else if comparison == 1 {
			newCoordinate = append(newCoordinate, prevCord[count])

			// Case 2.
			if count+1 == prevLen {
				newCoordinate = append(newCoordinate, 1) // Raised

				return newCoordinate

				// Case X.
			} else {
				newCoordinate = append(newCoordinate, prevCord[count+1]+1) // Inkrementera raise

				return newCoordinate

			}

		}
		count++

	}
	if count < nextLen {
		// Case 3. Ingen raised plats
		if nextCord[count] == 1 {
			newCoordinate = append(newCoordinate, 0) // Lower bound

			// varning: case saknas. om next också har nollor måste vi loopa?
			// WHILE LOOP.

			newCoordinate = append(newCoordinate, 1) // Lower bound

			// Case Y. Next längre än prev, men raised plats
		} else if nextCord[count] == 0 {
			// Case W

			for i := nextLen - count; i > 0; i-- {
				newCoordinate = append(newCoordinate, 0)

			}
			newCoordinate = append(newCoordinate, 1)

		} else {
			newCoordinate = append(newCoordinate, 1)

		}
	}

	return newCoordinate
}

func findPrevItem(insertionCoord []int, db LinkedList) *Item {
	prev := db.Head
	for prev.Next != nil {
		if CompareIndexes(prev.Next.Coordinate, insertionCoord) {
			break
		} else {
			prev = prev.Next
		}
	}
	return prev
}

func Insertion(letter string, coordinate []int, db LinkedList, uID int) LinkedList {

	prevItem := findPrevItem(coordinate, db)

	newItem := Item{Letter: letter, Coordinate: coordinate, ID: uID} //prev och next
	db.Length++

	// Case 4
	nextItem := prevItem.Next
	if nextItem == nil {
		db.Append(&newItem)
		return db
	}

	prevItem.Next = &newItem
	newItem.Prev = prevItem

	nextItem.Prev = &newItem
	newItem.Next = nextItem

	return db
}

func Deletion(coordinate []int, db LinkedList) LinkedList {
	prevItem := findPrevItem(coordinate, db)

	itemToRemove := prevItem.Next

	nextItem := itemToRemove.Next

	prevItem.Next = nextItem

	if nextItem != nil {
		nextItem.Prev = prevItem
	} else {
		db.Tail = prevItem
	}

	db.Length--

	return db
}

func (d *Document) Insert(letter string, uID int) {

	cursorPosCoordinate := d.CursorPosition.Coordinate

	// Case 4
	if d.CursorPosition.Next == nil {
		insertCoord := []int{cursorPosCoordinate[0] + 1}

		d.Textcontent = Insertion(letter, insertCoord, d.Textcontent, uID)

		d.CursorForward()

		return
	}
	cursorPosNextCoord := d.CursorPosition.Next.Coordinate

	insertCoord := getCoordinateProperties(cursorPosCoordinate, cursorPosNextCoord)

	d.Textcontent = Insertion(letter, insertCoord, d.Textcontent, uID)

	d.CursorForward()
}

func (d *Document) CursorForward() {
	if d.CursorPosition.Next != nil {
		d.CursorPosition = d.CursorPosition.Next
	}
}

func (d *Document) CursorBackwards() {
	// BOF Placeholder har ID 0
	if d.CursorPosition.Prev.ID != 0 {
		d.CursorPosition = d.CursorPosition.Prev

	}
}

func (d *Document) LoadInsert(index []int, uID int) {
	// tänkt att vara en funktion när andras ändringar uppdateras på din klient
	// hitta vart den ska in
	// lägg den där

}

func (d *Document) MoveCursor(index int) {
	docLength := d.Textcontent.Length

	if index > docLength || index < 0 {
		println("Error. Can't move cursor out of bounds")
	} else {

		// TODO: Loopa och gå igenom lista "x" antal gånger
		// Item som vi landar på är nya cursorposition
		var newPosition Item
		current := d.Textcontent.Head 
		for i := 0; i < index; i++ {
			current = current.Next 
		}
		newPosition = *current
		d.CursorPosition = &newPosition
	}
}

//OBS använder oss bara av current cursor position för deletion just nu
func (d *Document) Delete() {
	if d.CursorPosition.Prev != nil {
		savedCursor := d.CursorPosition

			// Link the previous node to the next node
			savedCursor.Prev.Next = savedCursor.Next
	
		if savedCursor.Next != nil {
			
			savedCursor.Next.Prev = savedCursor.Prev
		} else {
			// Om det är tailen
			d.Textcontent.Tail = savedCursor.Prev
		}

		//d.Textcontent = Deletion(d.CursorPosition.Coordinate, d.Textcontent)
		d.Textcontent.Length--

		if savedCursor.Prev.Next != nil{
		d.CursorPosition = savedCursor.Prev.Next
		} else{
			d.CursorPosition = d.Textcontent.Tail
		}
	}
}


func NewDocument() Document {

	// BOD = Beginning Of File
	BOF := Item{
		Letter:     "",
		Coordinate: []int{0},
		ID:         0,
		Next:       nil,
		Prev:       nil,
	}

	textContent := LinkedList{
		Head:   &BOF,
		Tail:   &BOF,
		Length: 0,
	}

	d := Document{
		Textcontent:    textContent,
		CursorPosition: &BOF,
	}

	return d
}
