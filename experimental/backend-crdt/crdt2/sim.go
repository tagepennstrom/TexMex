package main

import (
	"fmt"
)

func (d *Document) PrintDocument() {
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
		fmt.Print(current.Coordinate, current.Letter, "\n")
	}
	println("Result:", result)

}

func (d *Document) stringChecker() string{
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
	}
	return result
}

func (d *Document) cordsChecker() string{
	var result string
	for current := d.Textcontent.Head; current != nil; current = current.Next {
		result += fmt.Sprint(current.Coordinate) + ""
	}
	return result
}

func (d *Document) cleanDocument() {
	if d.Textcontent.Length > 0 {
		d.MoveCursor(1)	
	}
	for d.Textcontent.Length > 0 {
		d.Delete()
	}
}
func (d *Document) testCaseOne() int {
	d.cleanDocument()
	d.Insert("H", 1)
	d.Insert("E", 1)
	d.Insert("J", 1)
	d.Insert("!", 1)

	d.MoveCursor(3)
	d.Delete()
	d.MoveCursor(2)
	d.Insert("J", 1)

	if  d.stringChecker() == "HEJ!" && d.cordsChecker() == "[0][1][2][3][4]"{
		return 1;
	}else if d.stringChecker() != "HEJ!" && d.cordsChecker() != "[0][1][2][3][4]"{
		return 2;
	}else if d.stringChecker() != "HEJ!"{
		return 3;
	}else{
		return 4;
	}
}

func(d *Document) testCaseTwo() int {
	d.cleanDocument()
	d.Insert("A", 1)
	d.Insert("B", 1)
	d.Insert("D", 1)
	d.MoveCursor(2)
	d.Insert("C", 1)

	if  d.stringChecker() == "ABCD" && d.cordsChecker() == "[0][1][2][2 1][3]"{
		return 1;
	}else if d.stringChecker() != "ABCD" && d.cordsChecker() != "[0][1][2][2 1][3]"{
		return 2;
	}else if d.stringChecker() != "ABCD"{
		return 3;
	}else{
		return 4;
	}
}

func(d *Document) testCaseThree() int {
	d.cleanDocument()
	d.Insert("A", 1)
	d.Insert("B", 1)
	d.Insert("E", 1)
	d.MoveCursor(2)
	d.Insert("D", 1)
	d.MoveCursor(2)
	d.Insert("C", 1)

	if  d.stringChecker() == "ABCDE" && d.cordsChecker() == "[0][1][2][2 0 1][2 1][3]"{
		return 1;
	}else if d.stringChecker() != "ABCDE" && d.cordsChecker() != "[0][1][2][2 0 1][2 1][3]"{
		return 2;
	}else if d.stringChecker() != "ABCDE"{
		return 3;
	}else{
		return 4;
	}
}

func(d *Document) testCaseFour() int {
	d.cleanDocument()
	d.Insert("G", 1)
	d.Insert("U", 1)
	d.Insert("E", 1)
	d.Insert("N", 1)

	if  d.stringChecker() == "GUEN" && d.cordsChecker() == "[0][1][2][3][4]"{
		return 1;
	}else if d.stringChecker() != "GUEN" && d.cordsChecker() != "[0][1][2][3][4]"{
		return 2;
	}else if d.stringChecker() != "GUEN"{
		return 3;
	}else{
		return 4;
	}
}

func(d *Document) testCaseFive() int {
	d.cleanDocument()
	d.Insert("N", 1)
	d.MoveCursor(0)
	d.Insert("O", 1)
	d.MoveCursor(0)
	d.Insert("T", 1)
	d.MoveCursor(0)
	d.Insert("N", 1)
	d.MoveCursor(0)
	d.Insert("A", 1)


	if  d.stringChecker() == "ANTON" && d.cordsChecker() == "[0][0 0 0 0 1][0 0 0 1][0 0 1][0 1][1]"{
		return 1;
	}else if d.stringChecker() != "ANTON" && d.cordsChecker() != "[0][0 0 0 0 1][0 0 0 1][0 0 1][0 1][1]"{
		return 2;
	}else if d.stringChecker() != "ANTON"{
		return 3;
	}else{
		return 4;
	}
}

func(d *Document) testCaseX() int {
	d.cleanDocument()
	d.Insert("B", 1) //B[1]
	d.Insert("A", 1) //B[1]A[2]
	d.Insert("R", 1) //B[1]A[2]R[3]
	d.MoveCursor(2)  
	d.Insert("K", 1) //B[1]A[2]K[2 1]R[3]
	d.MoveCursor(2)
	d.Insert("C", 1) //B[1]A[2]C[2 0 1]K[2 1]R[3]
	d.MoveCursor(4)
	d.Insert("O", 1) //B[1]A[2]C[2 0 1]K[2 1]O[2,2]R[3]
	d.MoveCursor(4)
	d.Insert("D", 1) //B[1]A[2]C[2 0 1]K[2 1]D[2 1 1]O[2,2]R[3]
	d.MoveCursor(4)
	d.Insert("O", 1) //B[1]A[2]C[2 0 1]K[2 1]D[2 1 1]O[2 1 2]O[2,2]R[3]

	if  d.stringChecker() == "BACKDOOR" && d.cordsChecker() == "[0][1][2][2 0 1][2 1][2 1 1][2 1 2][2 2][3]"{
		return 1;
	}else if d.stringChecker() != "BACKDOOR" && d.cordsChecker() != "[0][1][2][2 0 1][2 1][2 1 1][2 1 2][2 2][3]"{
		return 2;
	}else if d.stringChecker() != "BACKDOOR"{
		return 3;
	}else{
		return 4;
	}
}

func main() {
	var doc Document = NewDocument()

	//doc.Insert("a", 1) // a [1]
	//doc.Insert("b", 1) // b [2]
	//doc.Insert("c", 1) // c [3]
	//doc.MoveCursor(2)
	//doc.Insert("d", 1) // d [2 1]

	// expected output fån print
	// [0]
	// [1] a
	// [2] b
	// [2 1] d
	// [3] c
	// Result: abdc

	//doc.PrintDocument()
	//for i := 1; i <= doc.Textcontent.Length; i++ {
	//	doc.MoveCursor(i)
	//	print("Cursor moved to position[", i, "]: ")
	//	fmt.Print("letter [", doc.CursorPosition.Letter, "] ", doc.CursorPosition.Coordinate, "\n")
	//}

	var failedTests string;
	var successTests string;


	println("TEST FOR CASE 1")
	println("---------------------------------------")

	switch doc.testCaseOne(){
		case 1:
			println("test case one was a success")
			successTests += "Test 1, "
		case 2:
			println("test case one was a fail neither the letters nor the coordinates were correct")
			println("Expected HEJ! but got: ", doc.stringChecker(), "Expected [0][1][2][3][4] but got:", doc.cordsChecker())
			failedTests += "Test 1, "
		case 3:
			println("test case one was a fail, the letters were not in the correct order")
			println("Expected HEJ! but got: ", doc.stringChecker())
			failedTests += "Test 1, "
		case 4:
			println("Test case one failed, the cords are not correct")
			println("Expected [0][1][2][3][4] but got:", doc.cordsChecker())
			failedTests += "Test 1, "
	}
	//println("string: ", doc.stringChecker(), "cords: ", doc.cordsChecker())
	println("---------------------------------------")
	println("\n")


	println("TEST FOR CASE 2")
	println("---------------------------------------")

	switch doc.testCaseTwo(){
	case 1:
		println("test case two was a success")
		successTests += "Test 2, "
	case 2:
		println("test case three was a fail neither the letters nor the coordinates were correct")
		println("Expected ABCD but got: ", doc.stringChecker(), "Expected [0][1][2][2 1][3] but got:", doc.cordsChecker())
		failedTests += "Test 2, "
	case 3:
		println("test case three was a fail, the letters were not in the correct order")
		println("Expected ABCD but got: ", doc.stringChecker())
		failedTests += "Test 2, "
	case 4:
		println("Test case three failed, the cords are not correct")
		println("Expected [0][1][2][2 1][3] but got:", doc.cordsChecker())
		failedTests += "Test 2, "
	}
	//println("string: ", doc.stringChecker(), "cords: ", doc.cordsChecker())
	println("---------------------------------------")
	println("\n")

	println("TEST FOR CASE 3")
	println("---------------------------------------")

	switch doc.testCaseThree(){
	case 1:
		println("test case two was a success")
		successTests += "Test 3, "
	case 2:
		println("test case three was a fail neither the letters nor the coordinates were correct")
		println("Expected ABCDE but got: ", doc.stringChecker(), "Expected [0][1][2][2 0 1][2 1][3] but got:", doc.cordsChecker())
		failedTests += "Test 3, "
	case 3:
		println("test case three was a fail, the letters were not in the correct order")
		println("Expected ABCDE but got: ", doc.stringChecker())
		failedTests += "Test 3, "
	case 4:
		println("Test case three failed, the cords are not correct")
		println("Expected [0][1][2][2 0 1][2 1][3] but got:", doc.cordsChecker())
		failedTests += "Test 3, "
	}
	//println("string: ", doc.stringChecker(), "cords: ", doc.cordsChecker())
	println("---------------------------------------")
	println("\n")


	println("TEST FOR CASE 4")
	println("---------------------------------------")
	switch doc.testCaseFour(){
	case 1:
		println("test case four was a success")
		successTests += "Test 4, "
	case 2:
		println("test case four was a fail neither the letters nor the coordinates were correct")
		println("Expected GUEN but got: ", doc.stringChecker(), "Expected [0][1][2][3][4] but got:", doc.cordsChecker())
		failedTests += "Test 4, "
	case 3:
		println("test case four was a fail, the letters were not in the correct order")
		println("Expected GUEN but got: ", doc.stringChecker())
		failedTests += "Test 4, "
	case 4:
		println("Test case four failed, the cords are not correct")
		println("Expected [0][1][2][3][4] but got:", doc.cordsChecker())
		failedTests += "Test 4, "
	}
	//println("string: ", doc.stringChecker(), "cords: ", doc.cordsChecker())
	println("---------------------------------------")
	println("\n")

	println("TEST FOR CASE 5")
	println("---------------------------------------")
	switch doc.testCaseFive(){
	case 1:
		println("test case five was a success")
		successTests += "Test 5, "
	case 2:
		println("test case five was a fail neither the letters nor the coordinates were correct")
		println("Expected ANTON but got: ", doc.stringChecker(), "Expected [0][0 0 0 0 1][0 0 0 1][0 0 1][0 1][1] but got:", doc.cordsChecker())
		failedTests += "Test 5, "
	case 3:
		println("test case five was a fail, the letters were not in the correct order")
		println("Expected ANTON but got: ", doc.stringChecker())
		failedTests += "Test 5, "
	case 4:
		println("Test case five failed, the cords are not correct")
		println("Expected [0][0 0 0 0 1][0 0 0 1][0 0 1][0 1][1] but got:", doc.cordsChecker())
		failedTests += "Test 5, "
	}
	//println("string: ", doc.stringChecker(), "cords: ", doc.cordsChecker())
	println("---------------------------------------")
	println("\n")

	println("TEST FOR CASE X")
	println("---------------------------------------")
	switch doc.testCaseX(){
	case 1:
		println("test case X was a success")
		successTests += "Test X, "
	case 2:
		println("test case X was a fail neither the letters nor the coordinates were correct")
		println("Expected BACKDOOR but got: ", doc.stringChecker(), "Expected [0][1][2][2 0 1][2 1][2 1 1][2 1 2][2 2][3] but got:", doc.cordsChecker())
		failedTests += "Test X, "
	case 3:
		println("test case X was a fail, the letters were not in the correct order")
		println("Expected BACKDOOR but got: ", doc.stringChecker())
		failedTests += "Test X, "
	case 4:
		println("Test case X failed, the cords are not correct")
		println("Expected [0][1][2][2 0 1][2 1][2 1 1][2 1 2][2 2][3] but got:", doc.cordsChecker())
		failedTests += "Test X, "
	}
	//println("string: ", doc.stringChecker(), "cords: ", doc.cordsChecker())
	println("---------------------------------------")
	println("\n")

	if len(successTests) >=2 {trimmed := successTests[:len(successTests)-2]; fmt.Println("Successfull tests:", trimmed +"!")}
	
	if len(failedTests) >=2 {trimmed := failedTests[:len(failedTests)-2]; 
			fmt.Println("failed tests:", trimmed +"!")
		}else{println(("No test failed"))}
}
