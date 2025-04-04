package main

import "fmt"

func (d *Document) PrintDocument() {
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
		fmt.Print(current.Coordinate, current.Letter, "\n")
	}
	println("Result:", result)

}

func main() {
	var doc Document = NewDocument()

	doc.Insert("a", 1) // a [1]
	doc.Insert("b", 1) // b [2]
	doc.Insert("c", 1) // c [3]

	// expected output fån print
	// [0]
	// [1] a
	// [2] b
	// [3] c
	// Result: abc

	doc.PrintDocument()
}
