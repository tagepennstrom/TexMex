package crdt

import (
	"os"
)

type Document struct {
	Head   *Item
	Tail   *Item
	Length int
}

type Item struct {
	Char   byte
	Location CoordT // istället för []int här
	ID       int
	Prev     *Item
	Next     *Item
}

type CoordT struct {
	Coordinate []int
	ID         int
}

func NewDocument() Document {
	var location CoordT = CoordT{
		Coordinate: []int{0},
		ID:         0,
	}

	// BOF = Beginning Of File
	BOF := Item{
		Char:   0,
		Location: location,
		ID:       0,
		Next:     nil,
		Prev:     nil,
	}
	return Document{
		Head:   &BOF,
		Tail:   &BOF,
		Length: 0,
	}
}

func DocumentFromStr(str string) Document {
	doc := NewDocument()
	for i := range len(str) {
		doc.Insert(str[i], i, 0)
	}
	return doc
}

func (doc *Document) ToString() string {
	str := ""
	item := doc.Head.Next
	for item != nil {
		str += string(item.Char)
		item = item.Next
	}
	return str
}

func (doc *Document) Insert(char byte, index int, uID int) {
	coordinate := doc.indexToCoordinate(index, uID)
	prevItem := findPrevItem(coordinate, *doc)

	newItem := Item{Char: char, Location: coordinate, ID: uID} //prev och next
	doc.Length++

	// Case 4
	nextItem := prevItem.Next
	if nextItem == nil {
		doc.append(&newItem)
	} else {
		prevItem.Next = &newItem
		newItem.Prev = prevItem

		nextItem.Prev = &newItem
		newItem.Next = nextItem
	}
}

func (doc *Document) indexToCoordinate(index int, uID int) CoordT {
	prevItem, caseFour := doc.indexToItem(index)

	// Case 4
	if caseFour {
		return getAppendCoordinate(prevItem.Location.Coordinate, uID)
	}

	nextItem := prevItem.Next
	coord := findIntermediateCoordinate(prevItem.Location, nextItem.Location)
	return CoordT{
		Coordinate: coord,
		ID:         uID,
	}
}

func (d *Document) indexToItem(index int) (Item, bool) {
	docLength := d.Length
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

		current := d.Head
		for range index {
			current = current.Next
		}

		newPosition = *current
	}
	return newPosition, atEnd
}

func getAppendCoordinate(prevCoord []int, uID int) CoordT {
	insertCoord := []int{prevCoord[0] + 1}
	var newLocation CoordT
	newLocation.Coordinate = insertCoord
	newLocation.ID = uID

	return newLocation
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

func (doc *Document) append(newItem *Item) {
	if doc.Tail == nil {
		println("Error: Detta bör aldrig hända")
		doc.Head = newItem
		doc.Tail = newItem
		return
	}

	doc.Tail.Next = newItem
	newItem.Prev = doc.Tail
	doc.Tail = newItem
	newItem.Next = nil
}

func findPrevItem(insertionCoord CoordT, doc Document) *Item {
	prev := doc.Head
	for prev.Next != nil {
		if compareIndexes(prev.Next.Location, insertionCoord) {
			break
		} else {
			prev = prev.Next
		}
	}
	return prev
}

// Returns true if c1 is smaller than c2
func compareIndexes(c1 CoordT, c2 CoordT) bool {

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
			println("Error: Coordinates can't have the same size and ID. This should not happen!")
			os.Exit(1)
		}

	}
	return len1 > len2
}
