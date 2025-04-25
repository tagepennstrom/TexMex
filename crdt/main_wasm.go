package crdt

import (
	"syscall/js"
)

var doc Document = NewDocument()

func insertWrap(this js.Value, args []js.Value) interface{} {

	if len(args) < 3 {
		println("Missing arguments")
		return nil
	}

	letter := args[0].String()

	index := args[1].Int()

	uID := args[2].Int()

	doc.LoadInsert(letter, index, uID)

	return nil
}

func cursorInsertWrap(this js.Value, args []js.Value) interface{} {

	if len(args) < 1 {
		println("Missing arguments")
		return nil
	}

	letter := args[0].String()

	uID := 69696969 // todo Obs ändra OBS

	doc.Insert(letter, uID)

	return nil
}

func registerCallbacks() {

	js.Global().Set("GoPrintDocument", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		doc.PrintDoc()
		return nil
	}))

	js.Global().Set("GoInsert", js.FuncOf(insertWrap))

	js.Global().Set("GoCursorInsert", js.FuncOf(cursorInsertWrap))

	println("Function callbacks registered")
}

func main() {
	//doc.LoadInsert("a", 1, 1)

	println("WASM is alive")
	registerCallbacks()
	select {} // keep running
}

/**

Var i CRDT directory innan denna i terminalen:
GOOS=js GOARCH=wasm go build -o ./frontend/static/wasm/main.wasm .

**/
