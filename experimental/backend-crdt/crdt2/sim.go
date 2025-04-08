package main

import (
	"fmt"
)

func (d *Document) Debug(extraInfo bool) {

	cursor := d.CursorPosition
	println("\n### Document Debug ###\n")

	println("# Cursor Info:")
	println("----------------------")
	fmt.Println("C coordinate:", cursor.Coordinate)
	fmt.Println("C letter:", cursor.Letter)

	if cursor.Prev != nil {
		fmt.Println("C previous:", cursor.Prev.Letter)
	} else {
		fmt.Println("C Prev is nil")
	}

	if cursor.Next != nil {
		fmt.Println("C next letter:", cursor.Next.Letter)
	} else {
		fmt.Println("C Next is nil")
	}
	println("----------------------")

	println("# Document Info:")
	println("----------------------")
	fmt.Println("Text content length:", d.Textcontent.Length)
	fmt.Println("TC Head letter:", d.Textcontent.Head.Letter)
	fmt.Println("TC Tail letter:", d.Textcontent.Tail.Letter)

	println("----------------------")

	println("# Textcontent:")
	println("----------------------")
	d.PrintDocument(extraInfo)
	println("----------------------")

}

func (d *Document) PrintDocument(verbose bool) {
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
		if verbose {
			fmt.Print("  ", current.Coordinate, current.Letter, "\n")
		}
	}
	println("Result:", result)

}

func (d *Document) CompileToText() string {
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
	}
	return result
}

func (d *Document) CoordToIndex(coord *Item) int {

	cur := d.Textcontent.Head
	index := 0
	for (cur.ID != coord.ID) || (cur.Letter != coord.Letter) {
		cur = cur.Next
		index++
	}
	return index
}

func (d *Document) DisplayWithCursor() {
	str := d.CompileToText()

	index := d.CoordToIndex(d.CursorPosition)

	result := str[:index] + "|" + str[index:]

	println(result)

}

func (d *Document) ExtractCoordinates() string {
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += fmt.Sprint(current.Coordinate) + ""
		fmt.Println("Loop: ", result)
	}
	return result
}

func (d *Document) CleanDocument() {
	if d.Textcontent.Tail.ID == 0 {
		println("Doc already cleared")
		return // doc already cleared
	}

	d.MoveCursor(d.Textcontent.Length)

	for d.Textcontent.Length > 0 {
		d.Delete()
	}
}

func main() {
	processInput()
	runAllTests() // <- Antons tester (ligger i tests.go)

	d := NewDocument()

	d.CleanDocument()
	d.Insert("A", 1)
	d.Insert("B", 1)
	d.Insert("D", 1)
	d.MoveCursor(2)
	d.Insert("C", 1)

	d.DisplayWithCursor()

	d.CleanDocument()

}
