package main

import (
	"syscall/js"

	"websocket-server/crdt"
)

var doc = crdt.NewDocument()

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

func updateDocument(document string, changes []crdt.Change, cursorIndex int) crdt.UpdatedDocMessage {
	doc = crdt.DocumentFromStr(document)
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

	return crdt.UpdatedDocMessage{
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

	changes := make([]crdt.Change, args[1].Length())
	for i := 0; i < len(changes); i++ {
		change := args[1].Index(i)
		changes[i] = crdt.Change{
			From: change.Get("from").Int(),
			To: change.Get("to").Int(),
			Text: change.Get("text").String(),
		}
	}
	res := updateDocument(document, changes, cursorIndex)
	var m = make(map[string]interface{})
	m["document"] = res.Document
	m["cursorIndex"] = res.CursorIndex

	return js.ValueOf(m)
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

Var i wasm directory innan denna i terminalen:
GOOS=js GOARCH=wasm go build -o ../frontend/static/wasm/main.wasm

**/
