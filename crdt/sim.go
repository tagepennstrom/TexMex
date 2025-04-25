package crdt

import (
	"fmt"
)

func (d *Document) Debug(extraInfo bool) {

	cursor := d.CursorPosition
	println("\n### Document Debug ###\n")

	println("# Cursor Info:")
	println("----------------------")
	fmt.Println("C coordinate:", cursor.Location)
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

	//println("THREE BIG BOOMS")
	//current := d.Textcontent.Head.Next
	//prev := d.Textcontent.Head
	//for i := 0; i < d.Textcontent.Length; i++{
	//fmt.Print("'"+current.Letter+"' " , "is number: " , i , "In the linked list")
	//print("\n")
	//fmt.Print("Is the prev my prev? ", prev.Letter, " ------------->>>> ", current.Prev.Letter)
	//print("\n")
	//d.moveCursorAndPrint()
	//prev = current
	//current = current.Next
	//}

}

func (d *Document) moveCursorAndPrint() {
	d.MoveCursor(1)

	println("Forward stepping")
	println("--------------------")
	for i := 0; i < d.Textcontent.Length; i++ {
		fmt.Print("cursor step", i, " = ")
		fmt.Print(d.CursorPosition.Letter)
		print("\n")
		d.CursorForward()
	}

	d.MoveCursor(d.Textcontent.Length)

	print("\n")

	println("Backward stepping")
	println("--------------------")
	for j := 0; j < d.Textcontent.Length; j++ {
		fmt.Print("cursor step", j, " = ")
		fmt.Print(d.CursorPosition.Letter)
		print("\n")
		d.CursorBackwards()
	}
}

func (d *Document) PrintDocument(verbose bool) {
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
		if verbose {
			fmt.Print("  ", current.Location, current.Letter, "\n")
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
	for cur != nil {
		if cur == coord {
			return index
		}
		cur = cur.Next
		index++
	}
	return -1 // Not found
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
		result += fmt.Sprint(current.Location) + ""
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
