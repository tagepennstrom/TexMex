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
	Letter   string
	Location CoordT // istället för []int här
	ID       int
	Prev     *Item
	Next     *Item
}

type CoordT struct {
	Coordinate []int
	ID         int
}

func (ll *LinkedList) Append(newItem *Item) {
	if ll.Tail == nil {
		println("Error: Detta bör aldrig hända")
		ll.Head = newItem
		ll.Tail = newItem
		return
	}

	ll.Tail.Next = newItem
	newItem.Prev = ll.Tail
	ll.Tail = newItem
	newItem.Next = nil
}

// Returns true if c1 is smaller than c2
func CompareIndexes(c1 CoordT, c2 CoordT) bool {

	coord1 := c1.Coordinate
	coord2 := c2.Coordinate
	len1 := len(coord1)
	len2 := len(coord2)

	count := 0

	for count < len1 && count < len2 {

		if coord1[count] > coord2[count] {
			return true //c1 biggest

		} else if coord1[count] < coord2[count] {
			return false //c2 biggest
		}
		count++
	}

	if len1 == len2 {

		if c1.ID < c2.ID {
			return true
		} else if c2.ID < c1.ID {
			return false
		} else {
			fmt.Errorf("Coordinates are identical")
			println("Error: Coordinates can't have the same size and ID. This should not happen!")
			os.Exit(1)
		}

	}
	return len1 > len2
}

func findIntermediateCoordinate(pCoord CoordT, nCoord CoordT) []int {

	prevCord := pCoord.Coordinate
	nextCord := nCoord.Coordinate

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

func findPrevItem(insertionCoord CoordT, db LinkedList) *Item {
	prev := db.Head
	for prev.Next != nil {
		if CompareIndexes(prev.Next.Location, insertionCoord) {
			break
		} else {
			prev = prev.Next
		}
	}
	return prev
}

func Insertion(letter string, coordinate CoordT, db LinkedList, uID int) LinkedList {

	prevItem := findPrevItem(coordinate, db)

	newItem := Item{Letter: letter, Location: coordinate, ID: uID} //prev och next
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

func Deletion(coordinate CoordT, db LinkedList) LinkedList {
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

func GetAppendCoordinate(prevCoord []int, uID int) CoordT {
	insertCoord := []int{prevCoord[0] + 1}
	var newLocation CoordT
	newLocation.Coordinate = insertCoord
	newLocation.ID = uID

	return newLocation
}

func (d *Document) Insert(letter string, uID int) {

	cursorPosCoordinate := d.CursorPosition.Location // TODO REWRITE, det här är oläsbart

	// Case 4
	if d.CursorPosition.Next == nil {
		location := GetAppendCoordinate(cursorPosCoordinate.Coordinate, uID)

		d.Textcontent = Insertion(letter, location, d.Textcontent, uID)
		d.CursorForward()

		return
	}
	cursorPosNextCoord := d.CursorPosition.Next.Location

	insertCoord := findIntermediateCoordinate(cursorPosCoordinate, cursorPosNextCoord)

	var location CoordT = CoordT{
		Coordinate: insertCoord,
		ID:         uID,
	}

	d.Textcontent = Insertion(letter, location, d.Textcontent, uID)
	d.CursorForward()

}

func (d *Document) CursorForward() {
	if d.CursorPosition.Next != nil {
		d.CursorPosition = d.CursorPosition.Next

	} else {
		println("Can't move cursor further.")
	}
}

func (d *Document) CursorBackwards() {
	// BOF Placeholder har ID 0
	if d.CursorPosition.ID != 0 {
		d.CursorPosition = d.CursorPosition.Prev

	} else {
		println("Error: Can't move cursor further back")
	}
}

func (d *Document) IndexToCoordinate(index int) (Item, bool) {
	docLength := d.Textcontent.Length
	var newPosition Item
	var atEnd bool = false

	if index >= docLength {
		index = docLength
		atEnd = true
	}

	if index < 0 {
		println("Error. Can't move cursor out of bounds")
		os.Exit(1)

	} else {

		current := d.Textcontent.Head
		for i := 0; i < index; i++ {
			current = current.Next
		}

		newPosition = *current
	}
	return newPosition, atEnd
}

func (d *Document) LoadInsert(letter string, index int, uID int) {

	prevItem, caseFour := d.IndexToCoordinate(index)
	var newCoordinate CoordT

	fmt.Println(prevItem.Location.Coordinate)

	// Case 4
	if caseFour {
		newCoordinate = GetAppendCoordinate(prevItem.Location.Coordinate, uID)
	}

	fmt.Println(newCoordinate)
	nextItem := prevItem.Next
	coord := findIntermediateCoordinate(prevItem.Location, nextItem.Location)

	var location CoordT = CoordT{
		Coordinate: coord,
		ID:         uID,
	}

	d.Textcontent = Insertion(letter, location, d.Textcontent, uID)
}

func (d *Document) MoveCursor(index int) {
	docLength := d.Textcontent.Length

	if index > docLength || index < 0 {
		println("Error. Can't move cursor out of bounds")
	} else {

		var newPosition Item
		current := d.Textcontent.Head
		for i := 0; i < index; i++ {
			current = current.Next
		}
		newPosition = *current
		d.CursorPosition = &newPosition
	}
}

// OBS använder oss bara av current cursor position för deletion just nu
func (d *Document) Delete() {
	if d.CursorPosition.Prev != nil {
		savedCursor := d.CursorPosition

		d.CursorBackwards()

		// Link the previous node to the next node
		savedCursor.Prev.Next = savedCursor.Next

		if savedCursor.Next != nil {

			savedCursor.Next.Prev = savedCursor.Prev
		} else {
			// Om det är tailen
			d.Textcontent.Tail = savedCursor.Prev
		}

		d.Textcontent.Length--

	}
}

func NewDocument() Document {

	var location CoordT = CoordT{
		Coordinate: []int{0},
		ID:         0,
	}

	// BOD = Beginning Of File
	BOF := Item{
		Letter:   "",
		Location: location,
		ID:       0,
		Next:     nil,
		Prev:     nil,
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
