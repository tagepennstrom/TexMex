package main

/*

import (
	"fmt"
	"os"
	"slices"
)

type Document struct {
	CursorPosition int
	Letters        []Item
}

type Item struct {
	Key    string
	Values []int
	ID     int
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
			// PRINTS - ("Case 1 or Z")

			//Case Z
			if count+1 < prevLen {
				newCoordinate = append(newCoordinate, prevCord[count+1]+1)
				// PRINTS - ("Confirmed as Case Z")

			}
			return newCoordinate

			// Case 2 and Case X. Ingen plats finns
		} else if comparison == 1 {
			newCoordinate = append(newCoordinate, prevCord[count])

			// Case 2.
			if count+1 == prevLen {
				newCoordinate = append(newCoordinate, 1) // Raised
				// PRINTS - ("Case 2")

				return newCoordinate

				// Case X.
			} else {
				newCoordinate = append(newCoordinate, prevCord[count+1]+1) // Inkrementera raise
				// PRINTS - ("Case X")

				return newCoordinate

			}

		}
		count++

	}
	if count < nextLen {
		// Case 3. Ingen raised plats
		if nextCord[count] == 1 {
			newCoordinate = append(newCoordinate, 0) // Lower bound
			// PRINTS - ("Case 3")

			// varning: case saknas. om next också har nollor måste vi loopa?
			// WHILE LOOP.

			newCoordinate = append(newCoordinate, 1) // Lower bound

			// Case Y. Next längre än prev, men raised plats
		} else if nextCord[count] == 0 {
			// Case W
			// PRINTS - ("Case W")

			for i := nextLen - count; i > 0; i-- {
				newCoordinate = append(newCoordinate, 0)

			}
			newCoordinate = append(newCoordinate, 1)

		} else {
			newCoordinate = append(newCoordinate, 1)
			// PRINTS - ("Case Y")

		}
	}

	return newCoordinate
}

func getIndexCoordinate(index int, db []Item) []int {
	return db[index].Values
}


func findInsertionCoordinate(index int, db []Item) []int {

	// Case 4. Insertion at end. [3, 1] EOF -> [4]
	if index == len(db)-1 {
		// PRINTS - ("Case 4")

		return []int{db[index].Values[0] + 1}

	}

	prev := getIndexCoordinate(index, db)
	next := getIndexCoordinate(index+1, db)

	// Case 5. Insertion at beginning. BOF [1] -> [0,1]
	if index == 0 {
		// PRINTS - ("Case 5")

		return getCoordinateProperties([]int{0}, next)
	}

	// Case 1,2,3,X,Y
	return getCoordinateProperties(prev, next)

}

func Insertion(letter string, index int, db []Item, uID int) []Item {
	coordinate := findInsertionCoordinate(index, db)

	item := Item{Key: letter, Values: coordinate, ID: uID}
	return slices.Insert(db, index+1, item)
}

func Deletion(index int, db []Item) []Item {
	if index < 0 || index >= len(db) {
		return db
	}
	return append(db[:index], db[index+1:]...)

}

func (d *Document) Insert(letter string, uID int) {

	d.Letters = Insertion(letter, d.CursorPosition, d.Letters, uID)
	d.CursorPosition++
}

func (d *Document) LoadInsert(index []int, uID int) {
	// hitta vart den ska in
	// lägg den där
	// updatera cursor om det behövs


}

func (d *Document) MoveCursor(index int) {
	if index >= len(d.Letters) || index < 0 {
		println("Error. Can't move cursor out of bounds")
	} else {
		d.CursorPosition = index
	}
}

func (d *Document) Delete() {
	if d.CursorPosition != 0 {
		d.Letters = Deletion(d.CursorPosition, d.Letters)
		d.CursorPosition--
	}
}

func NewDocument() Document {
	var doc Document
	doc.CursorPosition = 0
	doc.Letters = CRDTEmpty()
	return doc
}


func CompareIndexes(coord1 []int, coord2 []int) bool {

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
		fmt.Errorf("Coordinates have the same size")
		println("Error: Coordinates can't have the same size")
		os.Exit(1)
	}
	return len1 > len2
}


*/
