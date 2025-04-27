package crdt

import (
	"log"
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

func updateDocument(document string, changes []Change, cursorIndex int) UpdatedDocMessage {
	log.Println(document)
	doc = DocumentFromStr(document)
	doc.SetCursorAt(cursorIndex)
	for _, change := range changes {
		// TODO: give each user an ID
		uID := 1
		if change.Text == "" {
			for i := change.To; i > change.From; i-- {
				doc.DeleteAtIndex(i)
			}
		} else {
			for i := change.From; i <= change.To; i++ {
				doc.LoadInsert(string(change.Text[i - change.From]), change.From, uID)
			}
		}
	}

	return UpdatedDocMessage{
		Document: doc.ToString(),
		CursorIndex: doc.CursorIndex(),
	}
}

func updateDocumentWrap(this js.Value, args []js.Value) interface{} {

	if len(args) != 3 {
		println("Missing arguments")
		return nil
	}

	document := args[0].String()
	cursorIndex := args[2].Int()

	changes := make([]Change, args[1].Length())
	for i := 0; i < len(changes); i++ {
		change := args[1].Index(i)
		changes[i] = Change{
			From: change.Get("From").Int(),
			To: change.Get("To").Int(),
			Text: change.Get("Text").String(),
		}
	}
	return updateDocument(document, changes, cursorIndex)
}

func registerCallbacks() {

	js.Global().Set("GoPrintDocument", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		doc.PrintDoc()
		return nil
	}))

	js.Global().Set("GoInsert", js.FuncOf(insertWrap))

	js.Global().Set("GoCursorInsert", js.FuncOf(cursorInsertWrap))

	js.Global().Set("UpdateDocument", js.FuncOf(updateDocumentWrap))

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
